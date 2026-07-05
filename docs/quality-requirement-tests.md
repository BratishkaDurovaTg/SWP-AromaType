# Quality Requirement Tests

This document lists automated Quality Requirement Test evidence for the current
AromaType MVP v2 baseline. Manual demos and customer notes can support release
readiness, but they do not count as QRT evidence unless converted into an
automated check.

## QRT-001: Critical Module Coverage

**Linked quality requirement:** [QR-001](quality-requirements.md#qr-001-critical-module-testability)

**Verification method:** Automated CI coverage run.

**Environment:** GitHub Actions Ubuntu runner using the Go version from
`backend/go.mod`; latest local verification was run on 2026-07-03.

**Automated command or CI check:**

```bash
cd backend
go test -race -covermode=atomic -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Expected measurable result:** Critical backend packages stay at or above 30%
line coverage.

**Latest local evidence:**

- Total backend coverage: 36.3%
- `backend/internal/questionnaire`: 53.8%
- `backend/internal/http`: 52.5%
- `backend/internal/catalogbot`: 29.5%

**Evidence location:** Backend workflow logs and `backend-coverage` artifact.

## QRT-002: Recommendation Set Size

**Linked quality requirement:** [QR-002](quality-requirements.md#qr-002-recommendation-set-size)

**Verification method:** Automated unit test.

**Test data, setup, or environment:** Fixture answer weights and more than 5
candidate fragrances in `backend/internal/questionnaire/service_test.go`.

**Automated command or CI check:** `go test ./...`, specifically
`TestRecommendReturnsAtMostFiveItems`.

**Expected measurable result:** The recommendation response contains exactly 5
items when more than 5 catalog items match.

**Evidence location:** Backend workflow logs.

## QRT-003: Recommendation Determinism

**Linked quality requirement:** [QR-003](quality-requirements.md#qr-003-recommendation-determinism)

**Verification method:** Automated unit test.

**Test data, setup, or environment:** Fixture answer weights, tag names,
fragrance tags, and psychotype scores in
`backend/internal/questionnaire/service_test.go`.

**Automated command or CI check:** `go test ./...`, specifically
`TestRecommendRanksFragrancesAndBuildsProfile`.

**Expected measurable result:** The same fixture input produces the same top
fragrance, match percentage, profile bars, character traits, and key notes.

**Evidence location:** Backend workflow logs.

## QRT-004: Catalog Data Integrity

**Linked quality requirement:** [QR-004](quality-requirements.md#qr-004-catalog-data-integrity)

**Verification method:** Automated unit tests for catalog bot parsing helpers.

**Test data, setup, or environment:** Parser fixtures in
`backend/internal/catalogbot/parse_test.go`.

**Automated command or CI check:** `go test ./...`, specifically:

- `TestParseScores`
- `TestParseScoresRejectsInvalidValue`
- `TestParseVolumes`
- `TestValidateID`

**Expected measurable result:** Scores above 100 are rejected, volume strings are
parsed into structured volume/price pairs, and fragrance IDs are normalized to
safe slug values.

**Evidence location:** Backend workflow logs.

## QRT-005: Public API Contract Stability

**Linked quality requirement:** [QR-005](quality-requirements.md#qr-005-public-api-contract-stability)

**Verification method:** Automated HTTP router tests.

**Test data, setup, or environment:** In-memory test repository in
`backend/internal/http/router_test.go`.

**Automated command or CI check:** `go test ./...`, specifically:

- `TestHealthEndpointReturnsOK`
- `TestPublicQuestionnaireFlow`
- `TestAuthEndpointIsNotRegistered`

**Expected measurable result:** `/health` returns JSON status, the public
questionnaire/recommendation flow returns JSON recommendations, and removed JWT
registration returns `404 Not Found`.

**Evidence location:** Backend workflow logs.
