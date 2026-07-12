# ADR-002: Separate Telegram Catalog Bot

**Status:** Accepted

## Context

AromaType requires an admin interface for fragrance catalog management — adding,
editing, toggling, and updating photos of fragrances. Two approaches were
considered:

1. **Admin web UI** — add admin routes to the existing web frontend with
   authentication.
2. **Separate Telegram bot** — a dedicated Telegram bot for admin operations
   using the existing Telegram platform the team and customer already use.

The web frontend is a public-facing Mini App. Adding admin routes would require
building authentication, session management, and an admin UI — all of which are
outside the core user-facing product scope.

## Decision

Build a **separate Telegram catalog bot** as an independent Go binary sharing
the same database:

- Runs as a long-polling Telegram bot using the Bot API.
- Authenticates with a shared password (not per-user accounts).
- Supports commands: `/add`, `/edit`, `/list`, `/view`, `/toggle`, `/photo`,
  `/set`, `/cancel`, `/help`.
- Uses a 16-step wizard for new fragrance creation.
- Not started by the default `docker compose prod` command — requires explicit
  `--profile catalogbot` flag.

## Consequences and Tradeoffs

- **Positive:** No auth system needed in the public API — attack surface is
  reduced.
- **Positive:** Administrators use Telegram directly — no separate admin portal
  to build or maintain.
- **Positive:** Bot can be started or stopped independently from the API server.
- **Negative:** Requires a separate Telegram bot token and password management.
- **Negative:** Limited to predefined commands — no rich admin dashboard.
- **Negative:** Each admin needs the bot link and password; no per-user access
  control.

## Quality Requirements Addressed

- QR-004 (Fault tolerance) — catalog bot validates inputs before persisting;
  invalid data is rejected with clear error messages.
- QR-005 (Interoperability) — bot uses the same database and domain types as
  the API server.
