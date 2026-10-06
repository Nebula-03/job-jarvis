🤖 Job Jarvis — AI Job Intelligence Agent
An AI-powered job intelligence pipeline that transforms unstructured job-related emails into structured, trackable, and actionable information.

Status: 🚧 V1 — In Development
📌 Overview
Job Jarvis is a personal AI job-intelligence system designed to reduce the manual effort involved in managing job-search information.
It processes job-related emails from Gmail, extracts relevant information, analyzes job opportunities against a candidate profile, updates a job tracker, and surfaces important updates through notifications.
The system is designed around a human-in-the-loop approach.
AI interprets → Rules route → Human decides

Job Jarvis does not apply to jobs, reply to recruiters, submit applications, accept/decline offers, or make career decisions on behalf of the user.
🎯 Problem
Job searches generate a large amount of unstructured information:
- Application confirmations
- Job alerts
- Recruiter messages
- Interview invitations
- Assessments
- Rejections
- Status updates
- Offer notifications
Managing this information manually creates several problems:
- Important updates can be buried in an inbox.
- Application tracking becomes repetitive.
- Job descriptions require manual comparison against a candidate profile.
- Different email types require different processing paths.
- Relevant information is difficult to consolidate into one view.
Goal
Build an automated system that can:
Collect → Understand → Classify → Analyze → Track → Notify
while keeping final decisions with the user.
💡 Solution
Job Jarvis uses a hybrid architecture combining:
- n8n for workflow orchestration
- Groq for AI-based classification, extraction and analysis
- Gmail API for email ingestion and notifications
- Google Sheets for structured job tracking
- Docker for local infrastructure
Instead of using an LLM to control the entire workflow, AI is used for interpretation, while deterministic workflow logic controls execution.
🏗️ System Architecture
                    ┌─────────────────────┐
                    │      Gmail          │
                    │  Job-related Email  │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │        n8n           │
                    │ Filter + Fetch       │
                    │ Processed Check      │
                    │ Email Cleaning       │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    AI Reader        │
                    │       Groq          │
                    │ Classify + Extract  │
                    └──────────┬──────────┘
                               │
                               ▼
                    ┌─────────────────────┐
                    │    n8n Router       │
                    │ Deterministic Logic │
                    └──────┬───────┬──────┘
                           │       │
              ┌────────────┘       └────────────┐
              ▼                                 ▼
     ┌──────────────────┐             ┌──────────────────┐
     │ Google Sheets    │             │   AI Analyst     │
     │ Job Tracker      │             │      Groq         │
     └──────────────────┘             │   Job Fit        │
                                      └────────┬─────────┘
                                               │
                                               ▼
                                      ┌──────────────────┐
                                      │   Notifier       │
                                      │ Urgent + Digest  │
                                      └────────┬─────────┘
                                               │
                                               ▼
                                             Gmail

Note: The custom Go HTML-cleaning microservice is planned for Phase 2 and is not part of the current V1 architecture.

🔄 Workflow
1. Email Intake
n8n periodically checks the configured Gmail accounts for relevant job-related emails.
The initial workflow uses Gmail search filters to reduce unrelated email before AI processing.
2. Preprocessing
n8n removes unnecessary email content such as:
- HTML noise
- Signatures
- Quoted email history
- Tracking links
The goal is to provide the AI with the relevant visible content rather than the entire raw email.
3. AI Reader
The Reader Agent receives the cleaned email and determines:
- Whether it is job-related
- What type of event occurred
- Company
- Role
- Location
- Salary
- Required skills
- Application/job URL
- Interview or assessment information
- Recruiter information
- Other relevant fields
The output is returned as structured JSON.
4. Deterministic Routing
n8n evaluates the structured output and determines the next workflow path.
The router does not use an LLM.
This keeps workflow execution predictable and easier to debug.
5. Job Analysis
For job alerts and recruiter-related opportunities, the Analyst Agent evaluates the role against the candidate profile.
Possible categories:
- MATCH
- MATCH WITH LEARNING GAP
- LOW RELEVANCE
The system does not generate a numerical match score.
6. Tracking
Application-related information is written to the Google Sheets job tracker.
The system is designed to:
- Avoid duplicate jobs
- Avoid overwriting known values with blank/uncertain data
- Preserve uncertain information using CHECK: notes
- Maintain one row per job
7. Notifications
Important information is surfaced through:
- Urgent alerts
- Daily digest
The daily digest consolidates job-search updates into one message.
🧠 AI Agent Design
Job Jarvis does not depend on a single monolithic AI agent.
Instead, responsibilities are separated between specialized AI components.
AI Reader
Responsibility:
Understand the incoming email and convert it into structured information.

AI Analyst
Responsibility:
Analyze a job opportunity against the candidate profile.

AI Notifier
Responsibility:
Convert structured events and analysis into concise human-readable notifications.

This separation allows each AI component to have a smaller and more focused responsibility.
🛡️ Safety & Guardrails
Job Jarvis is intentionally designed with strict boundaries.
The system can:
✅ Read relevant emails
✅ Extract structured information
✅ Analyze job fit
✅ Update the job tracker
✅ Generate alerts
✅ Generate daily summaries  
The system cannot:
❌ Apply to jobs
❌ Submit applications
❌ Reply to recruiters
❌ Accept or decline offers
❌ Make career decisions
❌ Edit the user's resume
❌ Send emails to arbitrary recipients  
The user remains responsible for all external decisions and actions.
🛠️ Tech Stack
Technology	Purpose
n8n	Workflow orchestration
Groq	LLM inference
Gmail API	Email ingestion + notifications
Google Sheets API	Job tracking
Docker	Local infrastructure
Go	Planned Phase 2 email-cleaning microservice


📂 Project Structure
I'd keep this section aligned with what actually exists in the repo, rather than claiming directories that haven't been created yet.
For the planned structure:
job-jarvis/
│
├── services/
│   └── html-cleanser/        # Phase 2
│
├── workflows/                # n8n workflows
│
├── prompts/                  # AI prompts
│
├── docs/                     # Architecture & documentation
│
├── tests/                    # Test data / test cases
│
├── .env.example
├── .gitignore
├── docker-compose.yml
└── README.md

Don't create all these folders just because they're in the README. We'll update this once we actually structure the repo.
⚙️ Local Setup
Prerequisites
- Docker Desktop
- n8n
- Google Cloud project
- Gmail API enabled
- Google Sheets API enabled
- Groq API key
Start n8n
docker volume create n8n_data

docker run -d \
  --name n8n \
  --restart unless-stopped \
  -p 5678:5678 \
  -e GENERIC_TIMEZONE=Asia/Kolkata \
  -e TZ=Asia/Kolkata \
  -e EXECUTIONS_DATA_SAVE_ON_SUCCESS=none \
  -v n8n_data:/home/node/.n8n \
  docker.n8n.io/n8nio/n8n

Then open:
http://localhost:5678

🔐 Credentials & Security
Credentials should never be committed to GitHub.
Sensitive values are stored using n8n credentials/environment configuration.
Examples include:
Gmail OAuth credentials
Google Sheets OAuth credentials
Groq API key

.env files and credential files should be excluded through .gitignore.
🧪 Testing Strategy
Job Jarvis will be tested incrementally rather than deploying the complete workflow immediately.
Testing stages
Stage 1 — Dry Run
- No Google Sheets writes
- No outgoing emails
- Inspect AI outputs
Stage 2 — Controlled Testing
- Test against a copy of the tracker
- Send digest only to the development account
- Review classification accuracy
Stage 3 — Real Workflow
- Real tracker
- Both configured Gmail addresses
- Urgent notifications enabled
Test Dataset
The system will be tested against a mixture of:
- Application confirmations
- Job alerts
- Recruiter messages
- Interview invitations
- Assessments
- Rejections
- Newsletters
- Promotions
- Malformed emails
📊 Current V1 Scope
Implementing
- [ ] Gmail intake
- [ ] Email filtering
- [ ] Processed-email detection
- [ ] Email preprocessing
- [ ] AI Reader
- [ ] Structured JSON validation
- [ ] n8n routing
- [ ] AI Analyst
- [ ] Google Sheets tracker
- [ ] Urgent notifications
- [ ] Daily digest
- [ ] Error handling
- [ ] Test dataset
Already completed
- [x] Docker installed
- [x] Docker volume created
- [x] n8n running locally
- [x] n8n initial configuration
🚧 Known Limitations
The current V1 intentionally keeps the system simple.
Some architecture decisions remain under evaluation, including:
- Recruiter messages without a specific job
- Handling emails that require both tracking and urgent notification
- Public job-page processing
- Duplicate/merge behavior across multiple job sources
- AI failure handling
- Processed-email retention
These will be finalized during implementation and testing rather than assumed upfront.
🔮 Phase 2
After V1 has been used with real emails, a custom Go microservice may be introduced if email preprocessing becomes a bottleneck.
Go HTML Cleanser
n8n
  │
  │ POST /clean
  ▼
Go HTML Cleanser
  │
  │ cleaned text
  ▼
n8n → AI Reader

The service would programmatically process complex email HTML before AI inference.
This is intentionally not part of V1.
📈 Future Improvements
Potential future improvements include:
- Better email deduplication
- Improved failure recovery
- More robust job-source handling
- Additional structured storage
- Better observability
- Improved email preprocessing
- More sophisticated job analysis
- Migration from Google Sheets to a dedicated database if scale requires it
📝 Design Philosophy
Job Jarvis follows a few core principles:
1. Use AI where interpretation is required
LLMs handle classification, extraction and contextual analysis.
2. Use deterministic logic for execution
Workflow routing, validation and external actions remain controlled by explicit rules.
3. Keep the human in control
The system assists with information processing rather than making career decisions.
4. Start simple
V1 uses existing tools before introducing custom infrastructure.
5. Validate before persisting
AI output is validated before it can modify the job tracker.
🚀 Project Status
Current phase: V1 Development
The infrastructure is running locally, and the next stage is connecting Gmail and Google APIs before building the intake workflow.
