# ADR-003: Docker Compose Runtime

**Status:** Accepted

## Context

AromaType needs a consistent runtime environment across development machines
and the production VPS. The stack includes a Go backend, PostgreSQL database,
a Telegram catalog bot, and (in production) a Caddy reverse proxy. Options
considered:

1. **Docker Compose** — define all services in a single compose file.
2. **Manual host setup** — install Go runtime, PostgreSQL, Caddy directly on
   the host.
3. **Kubernetes** — full orchestration for a single-VPS deployment.

Manual setup introduces environment drift between team members. Kubernetes is
excessive for a single-node deployment.

## Decision

Use **Docker Compose** with two compose files:

- `docker-compose.yml` — local development with `profiles` to isolate the
  catalog bot.
- `docker-compose.prod.yml` — production deployment with Caddy, production
  environment, and explicit `--env-file` parameter.
- Multi-stage Docker build compiles Go binaries in a `golang:1.25-alpine`
  builder and copies them to a minimal `alpine:3.21` runtime image.

## Consequences and Tradeoffs

- **Positive:** Identical environment across all development machines and
  production — eliminates "works on my machine" issues.
- **Positive:** Single command to start the whole stack.
- **Positive:** Production and development share the same build and service
  definitions with environment-specific overrides.
- **Negative:** Docker overhead on resource-constrained machines.
- **Negative:** Build time for Go binary inside Docker on every stack start.
- **Negative:** Production secrets must be managed via `.env` file outside the
  compose configuration.

## Quality Requirements Addressed

- QR-001 (Testability) — consistent CI and local environment ensures tests pass
  identically in both contexts.
- QR-005 (Interoperability) — standardised compose network guarantees stable
  service discovery between components.
