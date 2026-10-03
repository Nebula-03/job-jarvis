package main

import (
	"encoding/json"
	"html"
	"log"
	"net/http"
	"strings"
	"unicode"

	xhtml "golang.org/x/net/html"
)

type CleanRequest struct {
	HTML    string `json:"html"`
	Subject string `json:"subject"`
	From    string `json:"from"`
}

type CleanResponse struct {
	SenderEmail        string `json:"sender_email"`
	EmailSubject       string `json:"email_subject"`
	CleanedTextContent string `json:"cleaned_text_content"`
	EstimatedTokens    int    `json:"estimated_tokens"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/clean", cleanHandler)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("HTML Cleanser running on http://localhost:8081")
	log.Println("POST /clean")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func cleanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
	defer r.Body.Close()

	var req CleanRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.HTML) == "" {
		http.Error(w, "html field is required", http.StatusBadRequest)
		return
	}

	cleanedText, err := cleanHTML(req.HTML)
	if err != nil {
		http.Error(w, "failed to parse HTML", http.StatusBadRequest)
		return
	}

	response := CleanResponse{
		SenderEmail:        strings.TrimSpace(req.From),
		EmailSubject:       strings.TrimSpace(req.Subject),
		CleanedTextContent: cleanedText,
		EstimatedTokens:    estimateTokens(cleanedText),
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func cleanHTML(rawHTML string) (string, error) {
	doc, err := xhtml.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	var walk func(*xhtml.Node)

	walk = func(node *xhtml.Node) {
		if node.Type == xhtml.ElementNode {
			switch strings.ToLower(node.Data) {
			case "script", "style", "noscript", "template":
				return
			}
		}

		if node.Type == xhtml.TextNode {
			text := html.UnescapeString(node.Data)

			if strings.TrimSpace(text) != "" {
				builder.WriteString(text)
				builder.WriteByte(' ')
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(doc)

	return normalizeWhitespace(builder.String()), nil
}

func normalizeWhitespace(text string) string {
	var builder strings.Builder
	spacePending := false

	for _, r := range text {
		if unicode.IsSpace(r) {
			spacePending = true
			continue
		}

		if spacePending && builder.Len() > 0 {
			builder.WriteByte(' ')
		}

		builder.WriteRune(r)
		spacePending = false
	}

	return strings.TrimSpace(builder.String())
}

func estimateTokens(text string) int {
	if strings.TrimSpace(text) == "" {
		return 0
	}

	// Approximation only.
	// Actual token count depends on the tokenizer used by the LLM.
	runeCount := len([]rune(text))

	return (runeCount + 3) / 4
}
