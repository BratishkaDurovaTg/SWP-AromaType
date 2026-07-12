# AGENTS.md — AromaType

This file provides operating guidance for coding agents working in the AromaType
repository. It complements `README.md` (project entry point) and `CONTRIBUTING.md`
(human contributor guide).

## Setup, Build, Test, and Verification Commands

### Full Stack (Docker Compose)

```bash
docker compose up --build
```

Opens at `http://localhost:5173` with backend at `http://localhost:8080`.

### Backend (Go)

```bash
cd backend

# Run server
go run ./cmd/api

# Run all tests with race detector and coverage
go test -race -covermode=atomic -coverprofile=coverage.out ./...

# View coverage report
go tool cover -func=coverage.out

# Check formatting (must pass — no output means clean)
gofmt -l .

# Auto-fix formatting
gofmt -w .

# Static analysis
go vet ./...

# Dependency vulnerability scan
govulncheck ./...
```

### Frontend (Vanilla JS)

```bash
cd frontend

# Install dependencies
npm install

# Run tests
npm test

# Dev server
python3 -m http.server 5173
```

## Repository Workflow

### Branching and PRs

1. Create a branch from the relevant issue page:
   `<issue-number>-short-description` (e.g., `42-add-questionnaire`).
2. Submit changes through a Pull Request (PR) to `dev` or `main`.
3. Every PR must:
   - Link to the relevant issue.
   - Be reviewed and approved by at least one other team member.
   - Pass all CI checks.
   - Complete the PR template checklist.
4. Use merge commits. Squash and rebase merging are disabled.
5. Do not push directly to `main` — it is branch-protected.

### CI Pipelines

All workflows run on PRs and pushes to `main`/`dev`:

| Workflow | What it checks |
|---|---|
| `Backend` | `gofmt`, `go vet`, `go test -race -covermode=atomic` with coverage artifact |
| `Frontend tests` | `npm ci`, `npm test` (Vitest) |
| `Additional QA` | `govulncheck ./...` (Go dependency vulnerability scan) |
| `Docker build` | Backend Docker image build |
| `Lychee` | Markdown link checking across the entire repo |

The latest `main` branch CI run must pass before submission.

### Issue Templates

Use the appropriate issue template:
- **User Story** — feature from the user's perspective with acceptance criteria
- **Feature Request** — product improvement or new capability
- **Bug Report** — problem description, reproduction steps, expected vs actual behaviour
- **Task** — technical, infrastructure, or documentation work

## Safety and Sensitive-Data Cautions

- **Never commit `.env` files, production secrets, or credentials.**
  `.gitignore` already ignores `.env`. Use `.env.example` as a template.
- **Never commit PII, real names, email addresses, or customer-identifying
  information** in public files. Use roles or pseudonyms (e.g., `customer`).
- **Never commit raw recordings, recording links, or exact private timecodes**
  to the repository. These are private-only artifacts.
- **Use only sanitized demo or test data** in public deployments, screenshots,
  API examples, and video demonstrations.
- If a secret is accidentally committed, follow the incident-response procedure:
  revoke the credential, make the repo private if needed, notify the TA, and
  remove the secret from history.

## Maintained Documentation

| Document | Purpose |
|---|---|
| `README.md` | Project entry point — description, access, links |
| `CONTRIBUTING.md` | Human contributor guide |
| `docs/customer-handover.md` | Customer handover status and transition scope |
| `docs/deployment.md` | Production deployment guide |
| `docs/testing.md` | Testing and CI status overview |
| `docs/quality-requirements.md` | ISO/IEC 25010 quality requirements with measurable scenarios |
| `docs/quality-requirement-tests.md` | Automated QRT definitions and evidence |
| `docs/user-acceptance-tests.md` | UAT scenarios for customer-facing workflows |
| `docs/user-stories.md` | User-story index with stable IDs |
| `docs/roadmap.md` | Sprint-by-Sprint delivery plan |
| `docs/definition-of-done.md` | Team's minimum completion standard |
| `docs/api/openapi.yaml` | OpenAPI specification |
| `docs/db-schema.md` | Database schema notes |
| `docs/catalog-bot.md` | Catalog bot commands and usage |

Update the relevant docs when the product, workflow, or configuration changes.
