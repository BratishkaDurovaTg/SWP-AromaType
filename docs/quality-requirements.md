# Quality Requirements

This document records maintained product quality requirements introduced for Assignment 4. Quality requirements use ISO/IEC 25010 sub-characteristics and are verified by automated quality requirement tests in [quality-requirement-tests.md](quality-requirement-tests.md).

## QR-001: Critical Module Testability

**ISO/IEC 25010 sub-characteristic:** Testability

**Scenario:** When a developer changes a critical backend product module under the standard CI environment, the module shall have automated tests that keep line coverage at or above 30%.

**Why this matters:** AromaType's recommendation and API modules support the core MVP user workflow. Defects in these modules can block the questionnaire, recommendations, or product details.

**Linked quality requirement tests:** [QRT-001](quality-requirement-tests.md#qrt-001-critical-module-coverage)

## QR-002: Recommendation Set Size

**ISO/IEC 25010 sub-characteristic:** Functional suitability

**Scenario:** When the recommendation service evaluates a catalog larger than the sample set size under the CI test environment, it shall return no more than 5 fragrance recommendations.

**Why this matters:** AromaType sells a compact sample set. The product must not overwhelm the user or promise more than the intended maximum number of samples.

**Linked quality requirement tests:** [QRT-002](quality-requirement-tests.md#qrt-002-recommendation-set-size)

## QR-003: Recommendation Determinism

**ISO/IEC 25010 sub-characteristic:** Functional correctness

**Scenario:** When the recommendation service receives the same set of answer option IDs under the CI test environment, it shall return the same ordered fragrance recommendations and match scores for the same fixture data.

**Why this matters:** The product does not promise magical AI guessing. Users and the team need a reproducible rule-based recommendation baseline before LLM explanations are added later.

**Linked quality requirement tests:** [QRT-003](quality-requirement-tests.md#qrt-003-recommendation-determinism)
