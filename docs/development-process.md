# Development Process — AromaType

This document describes the team's current actual development process as used
in the repository. It complements the [contributing guide](../CONTRIBUTING.md)
and the [agent guidance](../AGENTS.md).

## Backlog Management

### Product Backlog

The Product Backlog is managed through **GitHub Issues** and organised with a
**GitHub Projects board** (Kanban view). It acts as the single ordered source
of product work for the team.

Backlog items are classified as:

| Type | Description | Examples |
|---|---|---|
| **User Story** | Feature from the user's perspective with acceptance criteria | US-001: Complete Guided Questionnaire |
| **Technical PBI** | Implementation, refactoring, infrastructure, testing | TPBI-06: UI design fixing |
| **Bug** | Defect report with reproduction steps | — |
| **Course Task** | Course administration or reporting work (not a PBI) | Weekly report creation |

Backlog refinement is ongoing. Near-term, high-priority items are more detailed
than distant items.

### Sprint Backlog

Each Sprint has a **GitHub Milestone** and a **Sprint-specific view** on the
GitHub Projects board. The milestone records:

- Sprint Goal
- Start and finish dates
- Selected PBIs

A Sprint-selected PBI must be sufficiently ready to start — it must have:

- A clear expected outcome
- Description and context
- Acceptance criteria
- Story Point estimate (Modified Fibonacci: 1, 2, 3, 5, 8, 13, 20, 40, 100)
- An implementer
- A different reviewer

### Workflow States

Issues use the following Work Status values:

| Status | Meaning | Entry Criteria |
|---|---|---|
| `To Do` | Not yet ready to start | Product Backlog item with description |
| `Ready` | Selected for the Sprint, ready to begin | Sprint milestone assigned, estimated, implementer + reviewer assigned, acceptance criteria defined |
| `In Progress` | Work has started | Developer assigned the issue to themselves and started a branch |
| `Review` | PR/MR is open, review is active | PR/MR linked to the issue, all CI checks pass |
| `Done` | Completed | Acceptance criteria satisfied, Definition of Done met, PR merged |

## Git and Review Workflow

### Branching

```
main ────┬─────── merge commit ──────────── merge commit ────
         │                                      │
dev ─────┴─────── merge commit ────────────────┘
         │              │
feature  └──── branch ──┘
         <issue-number>-short-description
```

- `main` is the protected release branch.
- `dev` is the integration branch.
- Feature branches are created from the relevant issue page and named
  `<issue-number>-short-description` (e.g., `42-add-questionnaire`).

### Pull Requests

1. Developer creates a branch from the issue page.
2. Changes are committed on the feature branch.
3. A Pull Request is opened to `dev` or `main`.
4. The PR template must be completed:
   - Summary of changes
   - Related issue link
   - Testing performed checklist
   - Risk level
5. The PR links to the relevant issue. Automated dependency PRs are exempt.
6. The reviewer must be a different team member.
7. CI checks must pass before merge.
8. After approval, the PR is merged with a **merge commit** (squash and rebase
   are disabled).
9. The issue is closed manually after acceptance-criteria review.

### Review Requirements

- Every PR requires at least one approval from a team member who did not write
  the code.
- The reviewer verifies:
  - Acceptance criteria are satisfied
  - Tests pass and CI is green
  - Documentation is updated (OpenAPI, changelog, etc.)
  - Code follows project conventions
- The author may not approve their own PR.

## Configuration and Secrets Management

### Environment Variables

Configuration is supplied through environment variables. The repository contains:

- [`.env.example`](.env.example) — local development defaults (committed).
- [`.env.production.example`](.env.production.example) — production template
  with placeholders (committed).

Production `.env` files are created on the server and never committed to Git.

### Sensitive Data Rules

- `.env` files and `.env.production` files are ignored by `.gitignore`.
- Production secrets (database passwords, bot tokens, API keys) must never be
  committed.
- Use the `.env.example` template to document which variables are expected.
- If a secret is accidentally committed:
  1. Immediately revoke or rotate the exposed credential.
  2. Temporarily make the repository private.
  3. Notify the TA.
  4. Remove the secret from Git history.
  5. Document the incident privately.

### Sanitization

- Public documents use GitHub usernames, roles, or pseudonyms (e.g., `customer`)
  instead of real names.
- Screenshots and demo videos use only sanitised test data.
- Private recordings, private links, and exact timecodes are never committed to
  the public repository.

## Development Environment

### Reproducible Setup

The development environment is defined in `docker-compose.yml` and the backend
`Dockerfile`. All team members use the same containerised setup.

**Prerequisites:** Docker Desktop, Python 3, Git.

**Full stack start:**
```bash
docker compose up --build
cd frontend && python3 -m http.server 5173
```

**Backend only:**
```bash
cd backend && go run ./cmd/api
```

**Frontend only:** served by `python3 -m http.server 5173`, connects to the
running Docker backend at `localhost:8080`.

### Frontend Dependencies

```bash
cd frontend
npm install
npm test
```

Lockfile (`package-lock.json`) is committed for reproducible installs.

## CI and Automation

### CI Pipeline

Six GitHub Actions workflows run on every PR and push to `main`/`dev`:

| Workflow | Trigger | What it Runs | Required for Merge |
|---|---|---|---|
| `Backend` | PR, push to main/dev | `gofmt -l .`, `go vet ./...`, `go test -race -covermode=atomic` with coverage artifact | Yes |
| `Frontend tests` | PR, push to main/dev | `npm ci`, `npm test` (Vitest) | Yes |
| `Additional QA` | PR, push to main/dev | `govulncheck ./...` (dependency vulnerability scan) | Yes |
| `Docker build` | PR, push to main/dev | Docker image build for backend | Yes |
| `Lychee` | PR, push to main/dev | Markdown link checking across the repository | Yes |
| `Deploy backend` | Manual (`workflow_dispatch`) | SSH into VPS, `git pull`, `docker compose up -d --build` | No — triggered on demand |

### Quality Gates

The Definition of Done requires:

1. All acceptance criteria satisfied and verified.
2. PR reviewed and approved by another team member.
3. All CI checks pass (backend tests, frontend tests, formatting, vet,
   vulnerability scan, Docker build, link checking).
4. Automated tests implemented and passing for changed code.
5. Coverage for critical modules (`questionnaire`, `http`) stays at or above 30%.
6. Related documentation updated.
7. `CHANGELOG.md` updated for user-visible changes.
8. PR merged into the protected default branch.

### Deployment Automation

- **Development:** Manual `docker compose up --build`.
- **Production:** Manual deploy via GitHub Actions (`deploy-backend.yml`) or
  SSH with `docker compose -f docker-compose.prod.yml up -d --build`.
- **Catalog bot:** Not deployed by default — requires explicit profile flag.
- **No continuous deployment** — deployment is a manual decision.
