# Changelog

All notable changes to AromaType are documented in this file.

## [Unreleased]

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
