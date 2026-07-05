# ADR-0003: Use Docker Compose as the supported runtime topology

## Status

Accepted

## Context

The product uses a Go API, PostgreSQL, static frontend assets, an uploads
directory, and an optional catalog bot process. The team needs one repeatable
way to run the system locally and on a VPS.

## Decision

Use Docker Compose for the supported runtime topology. Local development runs
the backend and PostgreSQL from `docker-compose.yml`. Production deployment is
defined by `docker-compose.prod.yml` with Caddy serving static frontend assets
and reverse-proxying API traffic to the backend. The catalog bot is started via
an explicit Compose profile when admin catalog updates are needed.

## Consequences

- Backend, database, static frontend, uploads, and catalog bot wiring are
  reproducible.
- CI can build the backend Docker image before merge.
- Deployment can be repeated on a VPS without hand-running Go binaries.
- Secrets remain environment-specific and are not committed to the repository.

## Linked Quality Requirements

- [QR-001: Critical Module Testability](../../quality-requirements.md#qr-001-critical-module-testability)
- [QR-005: Public API Contract Stability](../../quality-requirements.md#qr-005-public-api-contract-stability)
