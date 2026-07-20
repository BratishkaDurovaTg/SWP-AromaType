# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]
No unreleased user-visible changes.

## [v1.4.0] - 2026-07-20

### Changed
- Removed recommendation category filters from the recommendations page.
- Improved cart scrolling so all items remain visible above the checkout button.
- Improved checkout validation by displaying an error message when required fields are empty.

### Fixed
- Removed unnecessary category filters displayed under "Ваша подборка".
- Fixed cart items appearing behind the "Оформить заказ" button.
- Fixed checkout validation by showing "Не все данные введены" when required fields are missing.

## [v1.3.0] - 2026-07-12

### Added

- Week 6 trial (handover-candidate) release for customer evaluation.
- Customer handover documentation (`docs/customer-handover.md`).
- Contributor guide (`CONTRIBUTING.md`) and agent guide (`AGENTS.md`).
- Updated customer-facing documentation to support independent product use.
- Transition-readiness documentation and customer trial materials.
- Hosted project documentation updates for the current trial release.

### Changed

- Updated the repository README to serve as the primary public entry point.
- Improved setup, access, and run instructions.
- Updated roadmap to reflect the remaining Week 6 and Week 7 work.
- Updated user stories for Sprint 4.
- Updated architecture, testing, quality requirements, and development process documentation.
- Improved customer-facing documentation based on Sprint 4 review.
- Refined the product based on customer feedback from MVP v2.

### Fixed

- Fixed documentation inconsistencies and outdated links.
- Fixed minor usability issues identified during customer review.
- Fixed issues discovered during customer trial and User Acceptance Testing.

---

## [v1.2.0] - 2026-07-05

### Added

- Repository baseline with protected `main` workflow.
- Go backend MVP with health check, Swagger UI, PostgreSQL, JWT auth, questionnaire API, and rule-based recommendations.
- Frontend MVP v0 smoke-check page.
- Lychee link-check workflow.
- Backend test workflow.
- Docker backend image build workflow.
- Manual backend deployment workflow for VPS.
- Issue templates and pull request template.
- Automated backend unit and integration tests, coverage reporting, quality requirement docs, and dependency vulnerability scan.
- Production Docker Compose and Caddy deployment configuration.
- Password-protected Telegram catalog bot for adding, viewing, editing, toggling, and uploading fragrance photos.
- Profile-specific result images for perfume types.
- Cart icon on the recommendation results screen.
- Add-to-cart button for recommendation cards.
- Frontend automated tests with Vitest and jsdom.
- GitHub Actions workflow for frontend tests.

### Changed

- Renamed Figma asset path from `ux:ui/v1.fig` to `ux-ui/v1.fig` for Windows compatibility.
- Updated the questionnaire to the 8-question psychotype structure from the latest product draft.
- Updated recommendation profile logic to score four psychotype tags: drive, focus, aesthetic, and power.
- Moved catalog management out of the public web app and into the separate Telegram bot workflow.
- Updated the perfume profile result screen according to customer feedback.
- Updated recommendation cards to show ordered numbers instead of a repeated `01`.
- Changed the sample set call-to-action text from `Заказать сет пробников` to `В корзину`.

### Removed

- Removed the `5 вариантов` label from the recommendation results header.
- Removed the `Доставка включена` text from the sample set block.
