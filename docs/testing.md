# Testing

This document is the canonical testing status artifact for AromaType. Assignment 4 checks introduced here are maintained product assets and should stay active unless replaced by documented equivalent or stronger checks.

## Critical Modules and Coverage

Coverage was measured locally with:

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

| Critical module | Why critical | Required line coverage | Current line coverage | Evidence |
|---|---|---:|---:|---|
| `backend/internal/questionnaire` | Core questionnaire matching, psychotype scoring, profile building, fragrance creation validation. | 30% | Measured by CI | `backend/internal/questionnaire/service_test.go`; backend CI coverage artifact |
| `backend/internal/http` | Public API routes, CORS, and upload validation helpers. | 30% | Measured by CI | `backend/internal/http/router_test.go`; `backend/internal/http/uploads_test.go`; backend CI coverage artifact |

Global backend coverage is measured by CI. It is lower than the tested service functions because database migration and repository code still depends on a live PostgreSQL environment and will be covered by stronger database integration tests in a later iteration.

## Automated Test Status

| Test type | Scope | Command or CI check | Latest result | Evidence |
|---|---|---|---|---|
| Unit tests | Recommendation ranking, psychotype scoring, max 5 recommendations, profile building, fragrance creation validation, upload content type mapping. | `go test ./...` | Passing locally | Test files under `backend/internal/**` |
| Integration tests | HTTP router with real questionnaire service and fake repositories for health, questions, recommendations, and removed auth endpoint behavior. | `go test ./...` | Passing locally | `backend/internal/http/router_test.go` |
| Automated QRTs | QR-001, QR-002, QR-003. | Backend workflow `test` job | Passing locally; CI runs on PRs and `main`/`dev` pushes | [quality-requirement-tests.md](quality-requirement-tests.md) |

## CI and QA Check Status

| Gate or check | Required for Done? | Latest protected-branch status | Evidence |
|---|---|---|---|
| Formatting | Yes | Runs on PRs and `main`/`dev` pushes | Backend workflow, `gofmt -l .` |
| Static analysis | Yes | Runs on PRs and `main`/`dev` pushes | Backend workflow, `go vet ./...` |
| Unit and integration tests | Yes | Runs on PRs and `main`/`dev` pushes | Backend workflow, `go test -race -covermode=atomic -coverprofile=coverage.out ./...` |
| Coverage reporting | Yes | Runs on PRs and `main`/`dev` pushes | Backend workflow coverage output and `backend-coverage` artifact |
| Docker build | Yes | Runs on PRs and `main`/`dev` pushes | Docker Build workflow |
| Link checking | Yes | Runs on PRs and `main`/`dev` pushes | Lychee workflow |
| Additional QA check | Yes | Runs on PRs and `main`/`dev` pushes | Additional QA workflow, `govulncheck ./...` |

## Additional QA Check Rationale

| QA objective or risk | Additional QA check | Scope | Latest result | Evidence | Limitations or follow-up |
|---|---|---|---|---|---|
| A vulnerable Go dependency or reachable standard-library vulnerability could expose product data, public API behavior, or deployment integrity. | Automated dependency vulnerability scan with `govulncheck`. | Go module dependencies and reachable backend code. | Passing locally; runs on PRs and `main`/`dev` pushes. | Additional QA workflow | It does not replace code review, secret scanning, or infrastructure hardening. Vulnerabilities may still require manual triage when upstream fixes are delayed. |

The team considered dependency vulnerability scanning, API contract checks, performance smoke tests, accessibility checks, and dependency freshness checks. `govulncheck` was selected first because the current backend exposes the public recommendation API and Go provides a reliable stack-native vulnerability scanner.

## Manual Evidence That Does Not Count as QRT

| Evidence | Scope | Result | Follow-up PBI or issue |
|---|---|---|---|
| Customer and team UI review | Questionnaire, profile, recommendation, and product screens. | Used as product feedback, not automated QRT evidence. | Track follow-up changes in GitHub Issues. |
| Swagger/manual API checks | Health, questionnaire, recommendation, and fragrance endpoints. | Useful smoke evidence while developing, not QRT unless automated. | Replace important manual checks with automated tests when a workflow becomes stable. |

## Maintained Gates After Assignment 4

The following gates remain active for later product work:

- Backend formatting check with `gofmt`.
- Backend static analysis with `go vet`.
- Backend unit and integration tests with race detector.
- Backend coverage reporting and critical-module coverage expectation of at least 30%.
- Docker backend image build.
- Lychee Markdown link checking.
- Additional QA dependency vulnerability scan with `govulncheck`.

Any future replacement must be documented here and must provide equivalent or stronger coverage of the same risk.
