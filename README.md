# 🤖 Job Jarvis

> AI-powered job intelligence pipeline for processing job-related Gmail,
> extracting structured information, analyzing job fit, tracking applications,
> and surfacing important updates.

🚧 **Status: V1 — In Development**

## Overview

Job Jarvis is a personal AI job-intelligence system designed to automate
the information-processing side of a job search while keeping all final
decisions with the user.

## Architecture

[architecture diagram]

## Core Design

AI interprets → Rules route → Human decides

## Tech Stack

- n8n — workflow orchestration
- Groq — LLM inference
- Gmail API — email ingestion and notifications
- Google Sheets API — job tracking
- Docker — local infrastructure
- Go — planned Phase 2 email-cleaning service

## Current Progress

- [x] Repository initialized
- [x] n8n running locally via Docker
- [ ] Gmail integration
- [ ] Google Sheets integration
- [ ] Groq integration
- [ ] Email intake workflow
- [ ] AI Reader
- [ ] Deterministic routing
- [ ] AI Analyst
- [ ] Notifications
- [ ] Validation and testing

## Roadmap

### V1
...
### Phase 2
Go HTML cleanser...
