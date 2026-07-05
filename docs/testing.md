# Testing

This document is the maintained testing status artifact for AromaType. It
records automated checks, quality requirement tests, and current known gaps for
the MVP v2 baseline.

## Automated Backend Checks

Backend checks are implemented in `.github/workflows/backend.yml` and run on
pull requests and pushes to `main` and `dev`.

| Check | Command | Purpose | Latest local result |
|---|---|---|---|
| Formatting | `test -z "$(gofmt -l .)"` | Prevent unformatted Go code. | Not rerun separately; included in CI. |
| Static analysis | `go vet ./...` | Catch common Go correctness issues. | Not rerun separately; included in CI. |
| Unit and integration tests | `go test -race -covermode=atomic -coverprofile=coverage.out ./...` | Validate recommendation logic, router behavior, upload helpers, and catalog parser helpers. | Passed on 2026-07-03. |
| Coverage report | `go tool cover -func=coverage.out` | Measure critical module coverage. | Passed on 2026-07-03. |

Latest local coverage evidence from 2026-07-03:

| Package | Coverage | Notes |
|---|---:|---|
| Total backend | 36.3% | Above the maintained 30% backend coverage threshold. |
| `backend/internal/questionnaire` | 53.8% | Critical recommendation and fragrance validation module. |
| `backend/internal/http` | 52.5% | Critical public API routing module. |
| `backend/internal/catalogbot` | 29.5% | Parser, validation, formatting, and keyboard helpers are covered; interactive Telegram flow still needs broader tests. |

## Automated Frontend Checks

Frontend tests are configured in `.github/workflows/frontend-tests.yml` and use
Vitest with jsdom. They cover the start screen, intro transition, public API
integration with mocked `fetch`, result rendering, profile image rendering,
recommendation card numbering, and cart button rendering.

Local note from 2026-07-03: `npm test` did not run on this machine because
frontend dev dependencies were not installed (`vitest: command not found`).
The local command sequence is:

```bash
cd frontend
npm install
npm test
```

## Other CI and QA Gates

| Workflow | Check | Purpose |
|---|---|---|
| `.github/workflows/docker-build.yml` | `docker build -f backend/Dockerfile -t aromatype-backend:ci .` | Confirms the backend Docker image can be built. |
| `.github/workflows/lychee.yml` | Lychee Markdown link check | Prevents broken links in repository Markdown files. |
| `.github/workflows/qa.yml` | `govulncheck ./...` | Detects reachable Go vulnerabilities in dependencies and standard library usage. |
| `.github/workflows/deploy-backend.yml` | Manual `workflow_dispatch` deploy | Deploys from `main` to a configured VPS when hosting is enabled. |

## Quality Requirement Test Matrix

| QRT | Requirement | Automated evidence |
|---|---|---|
| [QRT-001](quality-requirement-tests.md#qrt-001-critical-module-coverage) | [QR-001](quality-requirements.md#qr-001-critical-module-testability) | Backend coverage run and artifact. |
| [QRT-002](quality-requirement-tests.md#qrt-002-recommendation-set-size) | [QR-002](quality-requirements.md#qr-002-recommendation-set-size) | `TestRecommendReturnsAtMostFiveItems`. |
| [QRT-003](quality-requirement-tests.md#qrt-003-recommendation-determinism) | [QR-003](quality-requirements.md#qr-003-recommendation-determinism) | `TestRecommendRanksFragrancesAndBuildsProfile`. |
| [QRT-004](quality-requirement-tests.md#qrt-004-catalog-data-integrity) | [QR-004](quality-requirements.md#qr-004-catalog-data-integrity) | `backend/internal/catalogbot/parse_test.go`. |
| [QRT-005](quality-requirement-tests.md#qrt-005-public-api-contract-stability) | [QR-005](quality-requirements.md#qr-005-public-api-contract-stability) | `backend/internal/http/router_test.go`. |

## Manual Smoke Checks

Manual smoke checks are useful during demos but do not count as QRT evidence
until automated:

| Flow | Manual check |
|---|---|
| Mini App start | Open the frontend and verify the main screen loads. |
| Questionnaire | Complete all question steps and submit answers. |
| Results | Verify profile, recommendation list, and product cards render. |
| Product card | Open a fragrance, verify notes, accords, volume options, and image. |
| Catalog bot | Log in with the bot password, list products, edit a field, upload a photo, and verify product data through the public API. |

## Known Testing Gaps

- PostgreSQL repository and migration behavior are not covered by a real
  database integration test yet.
- Catalog bot conversation state, Telegram update handling, and photo download
  behavior need higher-level tests with a fake Telegram client.
- Frontend tests require installed npm dev dependencies locally.
- Hosted deployment smoke tests are currently disabled because the previous VPS
  deployment and DNS records were removed on 2026-07-03.

## Maintained Gates

The following checks must remain active or be replaced by equivalent stronger
checks:

- Go formatting.
- Go static analysis.
- Backend unit and integration tests with race detector.
- Critical backend module coverage at or above 30%.
- Docker backend image build.
- Lychee Markdown link check.
- Go vulnerability scan with `govulncheck`.
- Frontend Vitest/jsdom tests once dependencies are installed in the test
  environment.
