# Architecture

This directory documents the current AromaType architecture for MVP v2 planning.
The product is a Telegram Mini App/web frontend backed by a Go API, PostgreSQL,
and a separate Telegram catalog administration bot.

## System Responsibilities

| Component | Responsibility |
|---|---|
| Frontend Mini App | Presents the landing screen, questionnaire, profile result, recommendations, product cards, cart draft, and order draft UI. |
| Go API | Serves public questionnaire, recommendation, fragrance details, health, uploads, Swagger UI, and OpenAPI endpoints. |
| Recommendation service | Converts answer option IDs into weighted tags, psychotype scores, a profile, and up to 5 fragrance recommendations. |
| PostgreSQL | Stores questions, answer options, tags, fragrances, fragrance-tag weights, psychotype fields, volume options, and seeded MVP data. |
| Catalog Telegram bot | Provides password-protected catalog operations for admins outside the public user app. |
| GitHub Actions | Runs backend tests, Docker build, Lychee link checks, Go vulnerability scan, and frontend tests on PRs and protected branches. |

## Architecture Views

- [Static view](static-view.puml): modules and their dependencies.
- [Dynamic view](dynamic-view.puml): questionnaire-to-recommendation runtime flow.
- [Deployment view](deployment-view.puml): Docker-based runtime topology.

## Quality Requirement Links

- Recommendation logic supports [QR-002](../quality-requirements.md#qr-002-recommendation-set-size)
  and [QR-003](../quality-requirements.md#qr-003-recommendation-determinism).
- Catalog parsing and validation supports [QR-004](../quality-requirements.md#qr-004-catalog-data-integrity).
- Public API smoke behavior supports [QR-005](../quality-requirements.md#qr-005-public-api-contract-stability).

## Current Constraints

- Public user authentication is intentionally removed for the MVP flow.
- Product catalog administration is handled by a separate Telegram bot with a
  password gate, not by the public Mini App.
- The recommendation engine is rule-based and deterministic; it does not use an
  LLM for scoring.
- Deployment is Docker based. The previous VPS deployment was intentionally
  removed on 2026-07-03, so the deployment diagram describes the supported
  topology rather than an active production instance.
