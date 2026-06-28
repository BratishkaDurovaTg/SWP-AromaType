# Quality Requirement Tests

This document lists automated QRT evidence. Manual reviews, customer notes, and UAT observations can support release readiness, but they do not count as QRT evidence unless a later assignment explicitly allows it.

## QRT-001: Critical Module Coverage

**Linked quality requirement:** [QR-001](quality-requirements.md#qr-001-critical-module-testability)

**Verification method:** Automated CI coverage run.

**Test data, setup, or environment:** GitHub Actions Ubuntu runner using the Go version from `backend/go.mod`.

**Automated command or CI check:** Backend workflow job `test`, command `go test -race -covermode=atomic -coverprofile=coverage.out ./...`, followed by `go tool cover -func=coverage.out`.

**Expected measurable result:** Critical backend packages remain at or above 30% line coverage.

**Evidence location:** Coverage artifact `backend-coverage` and backend workflow logs.

## QRT-002: Recommendation Set Size

**Linked quality requirement:** [QR-002](quality-requirements.md#qr-002-recommendation-set-size)

**Verification method:** Automated unit test.

**Test data, setup, or environment:** Fixture answer weights and 8 candidate fragrances in `backend/internal/questionnaire/service_test.go`.

**Automated command or CI check:** `go test ./...`, specifically `TestRecommendReturnsAtMostFiveItems`.

**Expected measurable result:** The recommendation response contains exactly 5 items when more than 5 catalog items match.

**Evidence location:** Backend workflow logs.

## QRT-003: Recommendation Determinism

**Linked quality requirement:** [QR-003](quality-requirements.md#qr-003-recommendation-determinism)

**Verification method:** Automated unit test.

**Test data, setup, or environment:** Fixture answer weights, tags, and fragrances in `backend/internal/questionnaire/service_test.go`.

**Automated command or CI check:** `go test ./...`, specifically `TestRecommendRanksFragrancesAndBuildsProfile`.

**Expected measurable result:** The highest scoring fragrance is first, the top match percent is `99`, and the psychotype profile is built from the same weighted tag inputs.

**Evidence location:** Backend workflow logs.
