# Architecture — AromaType

This document describes the current delivered architecture of AromaType, a
Telegram Mini App for personalised fragrance discovery. It complements the
[README.md](../../README.md) and the [deployment guide](../deployment.md).

## Overview

AromaType follows a **modular monolith** pattern with clear package boundaries
and two separate entry points sharing a common database.

```
Telegram Mini App (Vanilla JS SPA) ──HTTP──> Go API Server ──> PostgreSQL 16
                                                   │
                                        Catalog Bot (Go) ──> PostgreSQL 16
```

The front end is a static single-page application embedded in Telegram via the
Telegram WebApp SDK. The back end provides a REST API for the questionnaire,
fragrance recommendations, and product catalog. A separate Telegram bot binary
handles admin catalog management.

## Architecture Views

- [Static view](static-view/) — Component structure, packages, and module boundaries
- [Dynamic view](dynamic-view/) — Runtime interactions, sequence flows, and data paths
- [Deployment view](deployment-view/) — Infrastructure, network topology, and hosting

## Architecture Decision Records

The following ADRs document key architectural decisions:

| ADR | Title | Status | Quality Requirements |
|---|---|---|---|
| [ADR-001](adr/ADR-001-rule-based-recommendation-engine.md) | Rule-Based Recommendation Engine | Accepted | QR-002, QR-003 |
| [ADR-002](adr/ADR-002-separate-catalog-bot.md) | Separate Telegram Catalog Bot | Accepted | QR-004, QR-005 |
| [ADR-003](adr/ADR-003-docker-compose-runtime.md) | Docker Compose Runtime for Development and Production | Accepted | QR-001, QR-005 |

### ADR Index

Each ADR is a standalone markdown file in [`docs/architecture/adr/`](adr/)
following the same structure:
- **Context** — why a decision was needed
- **Decision** — what was chosen
- **Consequences and tradeoffs** — what the decision affects
- **Quality requirements addressed** — which QR scenarios are supported

## Quality Requirements Addressed

| QR | Sub-characteristic | Architecture Support |
|---|---|---|
| QR-001 | Testability | Modular package structure enables focused unit testing of the recommendation engine in isolation. |
| QR-002 | Functional appropriateness | Rule-based design guarantees at most 5 recommendations per request. |
| QR-003 | Functional correctness | Deterministic scoring produces identical results for identical inputs. |
| QR-004 | Fault tolerance | Catalog bot validates all inputs before persisting; API rejects invalid recommendation requests. |
| QR-005 | Interoperability | Stateless JSON API with no auth requirement simplifies front-end integration. |

## Key Technology Choices

| Layer | Technology | Rationale |
|---|---|---|
| Front end | Vanilla JavaScript, HTML5, CSS | Lightweight SPA with no framework overhead; Telegram Mini App SDK dependency. |
| Back end | Go 1.25 | Fast startup, small binary, strong standard library mux (Go 1.22+), good concurrency. |
| Database | PostgreSQL 16 with JSONB | Flexible fragrance metadata (notes, scores, volumes) without schema migrations for every field. |
| Reverse proxy | Caddy 2 | Automatic HTTPS, easy configuration, built-in SPA fallback. |
| Runtime | Docker Compose | Consistent environment across development and production. |
| CI/CD | GitHub Actions | 6 workflows covering format, lint, test, build, deploy, and link checking. |
