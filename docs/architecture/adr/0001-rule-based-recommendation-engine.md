# ADR-0001: Use a deterministic rule-based recommendation engine

## Status

Accepted

## Context

AromaType must recommend fragrances from a prepared catalog after a short
questionnaire. The product does not promise that AI magically guesses a perfume.
The current MVP requires explainable matching, repeatable results, and a maximum
of 5 recommended samples.

## Decision

Use a deterministic rule-based recommendation engine in Go. Each answer option
maps to weighted tags and psychotype dimensions. Each fragrance stores tags,
psychotype, psychotype scores, accords, notes, volume options, and active status.
The service scores fragrance matches, sorts them deterministically by score and
name, builds a user profile, and returns no more than 5 items.

## Consequences

- Recommendations are reproducible for tests, demos, and customer review.
- The team can tune tags and weights without training an ML model.
- The engine is simple enough to validate with unit tests.
- LLM-based explanations can be added later without replacing the scoring core.

## Linked Quality Requirements

- [QR-002: Recommendation Set Size](../../quality-requirements.md#qr-002-recommendation-set-size)
- [QR-003: Recommendation Determinism](../../quality-requirements.md#qr-003-recommendation-determinism)
