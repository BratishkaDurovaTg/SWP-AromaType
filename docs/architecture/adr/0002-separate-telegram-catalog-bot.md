# ADR-0002: Manage the catalog through a separate Telegram admin bot

## Status

Accepted

## Context

The public MVP flow does not require user registration. The team still needs a
way to add and update fragrances, photos, notes, accords, psychotypes, volume
options, and active status without exposing admin controls to ordinary users.

## Decision

Keep catalog administration outside the public Mini App and implement it as a
separate password-protected Telegram bot. The bot uses the same PostgreSQL
database as the public API, stores uploaded product photos in the shared uploads
volume, and validates catalog fields before saving them.

## Consequences

- The user-facing Mini App remains focused on discovery and ordering.
- Admin access is isolated from public frontend routes.
- Catalog changes can be performed from Telegram without building a full admin
  web interface for MVP v2.
- Bot validation becomes part of the quality strategy because invalid catalog
  data can break recommendations and product cards.

## Linked Quality Requirements

- [QR-004: Catalog Data Integrity](../../quality-requirements.md#qr-004-catalog-data-integrity)
- [QR-005: Public API Contract Stability](../../quality-requirements.md#qr-005-public-api-contract-stability)
