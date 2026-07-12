# Contributing to AromaType

## Getting Started

1. Clone the repository.
2. Run the full stack with Docker Compose:
   ```bash
   docker compose up --build
   ```
3. Open `http://localhost:5173` (frontend) and `http://localhost:8080/health` (backend).

For detailed setup instructions, see `README.md`.

## Before Submitting Changes

- Run backend tests:
  ```bash
  cd backend
  go test -race -covermode=atomic -coverprofile=coverage.out ./...
  ```
- Check formatting (no output means clean):
  ```bash
  gofmt -l .
  ```
- Run static analysis:
  ```bash
  go vet ./...
  ```
- Run dependency vulnerability scan:
  ```bash
  govulncheck ./...
  ```
- Run frontend tests:
  ```bash
  cd frontend
  npm test
  ```
- Verify `docker compose up --build` starts without errors.
- If you changed the API, update `docs/api/openapi.yaml`.
- Add a changelog entry in `CHANGELOG.md` under `[Unreleased]` for user-visible changes.

## Branching and Pull Requests

1. Create a branch from the relevant GitHub issue:
   `<issue-number>-short-description` (e.g., `42-add-questionnaire`).
2. Make your changes on that branch.
3. Open a Pull Request (PR) to `dev` or `main`.
4. Complete the PR template checklist.
5. Request review from at least one other team member.
6. After approval, merge with a merge commit (squash and rebase are disabled).
7. Link the related issue in the PR description.

## Review Expectations

- Every PR must be approved by at least one team member who did not author the changes.
- The reviewer verifies:
  - Acceptance criteria are satisfied
  - Tests pass and CI checks are green
  - Documentation is updated if needed
  - Changelog entry is present for user-visible changes
- The author may not approve their own PR.

## Documentation

Key documents for contributors:

| Document | Purpose |
|---|---|
| `README.md` | Project overview, setup, and access |
| `AGENTS.md` | Agent-specific setup and workflow guidance |
| `docs/deployment.md` | Production deployment and operations |
| `docs/testing.md` | CI gates, test coverage, QA checks |
| `docs/quality-requirements.md` | Quality requirements and scenarios |
| `docs/quality-requirement-tests.md` | Automated QRT definitions and evidence |
| `docs/api/openapi.yaml` | API contract specification |
| `docs/db-schema.md` | Database schema notes |
| `docs/definition-of-done.md` | Completion standard for merged work |
| `docs/customer-handover.md` | Customer transition scope and status |
| `docs/catalog-bot.md` | Catalog admin bot commands |

Update the relevant documentation when your changes affect the product, workflow, or configuration.
