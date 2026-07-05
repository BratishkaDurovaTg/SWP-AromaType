# Quality Requirements

This document records maintained AromaType quality requirements for the current
MVP v2 baseline. Requirements use ISO/IEC 25010 quality characteristics and are
verified by automated quality requirement tests in
[quality-requirement-tests.md](quality-requirement-tests.md).

## QR-001: Critical Module Testability

**ISO/IEC 25010 characteristic:** Maintainability

**ISO/IEC 25010 sub-characteristic:** Testability

**Scenario:** When a developer changes a critical backend product module under
the standard CI environment, the module shall keep automated line coverage at or
above 30%.

**Measurement:** `go test -race -covermode=atomic -coverprofile=coverage.out ./...`
and `go tool cover -func=coverage.out`.

**Current coverage evidence from 2026-07-03:**

- Total backend coverage: 36.3%.
- `backend/internal/questionnaire`: 53.8% line coverage in the latest local run.
- `backend/internal/http`: 52.5% line coverage in the latest local run.
- `backend/internal/catalogbot`: 29.5% line coverage in the latest local run.

**Why this matters:** These modules support the public questionnaire,
recommendations, product details, and API routing. Defects here can block the
main MVP flow.

**Linked quality requirement tests:** [QRT-001](quality-requirement-tests.md#qrt-001-critical-module-coverage)

**Linked ADRs:** [ADR-0003](architecture/adr/0003-docker-compose-runtime.md)

## QR-002: Recommendation Set Size

**ISO/IEC 25010 characteristic:** Functional suitability

**ISO/IEC 25010 sub-characteristic:** Functional appropriateness

**Scenario:** When the recommendation service evaluates a catalog larger than
the sample set size, it shall return no more than 5 fragrance recommendations.

**Measurement:** Automated unit test with more than 5 matching fragrance
candidates.

**Why this matters:** AromaType recommends a compact set of samples. The product
must not overwhelm the user or promise more than the intended maximum sample
count.

**Linked quality requirement tests:** [QRT-002](quality-requirement-tests.md#qrt-002-recommendation-set-size)

**Linked ADRs:** [ADR-0001](architecture/adr/0001-rule-based-recommendation-engine.md)

## QR-003: Recommendation Determinism

**ISO/IEC 25010 characteristic:** Functional suitability

**ISO/IEC 25010 sub-characteristic:** Functional correctness

**Scenario:** When the recommendation service receives the same answer option
IDs and fixture catalog data, it shall return the same ordered fragrance
recommendations, match percentages, and profile values.

**Measurement:** Automated unit test using fixed answer weights, fragrance tags,
and psychotype scores.

**Why this matters:** The product uses honest rule-based matching, not hidden
prediction. Reproducible recommendations are required for debugging, tuning, and
customer review.

**Linked quality requirement tests:** [QRT-003](quality-requirement-tests.md#qrt-003-recommendation-determinism)

**Linked ADRs:** [ADR-0001](architecture/adr/0001-rule-based-recommendation-engine.md)

## QR-004: Catalog Data Integrity

**ISO/IEC 25010 characteristic:** Reliability

**ISO/IEC 25010 sub-characteristic:** Fault tolerance

**Scenario:** When an admin enters catalog data through the Telegram catalog bot,
the backend shall reject invalid parser inputs for structured fields before they
can be saved as fragrance data.

**Measurement:** Automated parser tests for psychotype score ranges, volume
format, and product ID normalization.

**Why this matters:** The public app depends on catalog data for product cards
and recommendations. Invalid score, volume, or ID formats can break matching or
render incorrect product information.

**Linked quality requirement tests:** [QRT-004](quality-requirement-tests.md#qrt-004-catalog-data-integrity)

**Linked ADRs:** [ADR-0002](architecture/adr/0002-separate-telegram-catalog-bot.md)

## QR-005: Public API Contract Stability

**ISO/IEC 25010 characteristic:** Compatibility

**ISO/IEC 25010 sub-characteristic:** Interoperability

**Scenario:** When the frontend calls the public MVP API, the backend shall keep
the questionnaire and recommendation endpoints available as JSON endpoints and
shall not expose removed JWT registration endpoints.

**Measurement:** Automated HTTP router tests for `/health`,
`GET /api/questions`, `POST /api/recommendations`, and the removed
`POST /api/auth/register` route.

**Why this matters:** The current product flow has no public registration. The
frontend depends on stable anonymous questionnaire and recommendation endpoints.

**Linked quality requirement tests:** [QRT-005](quality-requirement-tests.md#qrt-005-public-api-contract-stability)

**Linked ADRs:** [ADR-0002](architecture/adr/0002-separate-telegram-catalog-bot.md),
[ADR-0003](architecture/adr/0003-docker-compose-runtime.md)
