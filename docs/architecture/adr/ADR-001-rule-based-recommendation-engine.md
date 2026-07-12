# ADR-001: Rule-Based Recommendation Engine

**Status:** Accepted

## Context

AromaType needs to recommend fragrances based on user questionnaire responses.
The team evaluated two approaches:

1. **Rule-based scoring** — deterministic matching using psychotype tag weights
   and dot-product similarity.
2. **ML/Learned model** — train a model on user behaviour data to predict
   preferences.

At MVP stage there is no historical user data to train a model. The product
needs to work from launch with seed data. The recommendation logic must be
explainable to the customer and auditable by the team.

## Decision

Use a **rule-based recommendation engine** with deterministic psychotype scoring:

- Each answer option maps to psychotype tags (drive, focus, aesthetic, power)
  with configurable weights.
- Each fragrance defines its psychotype scores.
- The match score is the dot product of user and fragrance psychotype vectors.
- Tag-level matching adds additional weighted scores for non-psychotype tags.
- Results are sorted by total score, limited to 5, and normalised to a 70-99%
  match percentage.

## Consequences and Tradeoffs

- **Positive:** Deterministic results — same inputs always produce same outputs.
- **Positive:** Works from day one with seed data — no training data needed.
- **Positive:** Easy to debug and explain to the customer.
- **Positive:** Simple to unit test with fake repositories.
- **Negative:** Cannot learn from user behaviour over time without manual rule
  adjustments.
- **Negative:** Rule quality depends on how well psychotype tags match real
  user preferences.

## Quality Requirements Addressed

- QR-002 (Functional appropriateness) — guarantees at most 5 recommendations.
- QR-003 (Functional correctness) — deterministic, auditable scoring.
