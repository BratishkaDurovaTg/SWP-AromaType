# Customer Handover — AromaType

## Current Product Status

AromaType is a Telegram Mini App for personalised fragrance discovery through
style, situations, and feelings rather than perfume terminology. Users complete
a guided questionnaire, receive up to 5 recommended fragrances with explanations,
and order sample sets through a configured Telegram contact.

- **Current version:** MVP v3 (v1.4.0)
- **Handover level:** Ready for independent use
- **Customer-confirmation status:** Accepted with follow-up items

## How the Customer Accesses and Uses the Product

| What | How |
|---|---|
| Telegram Mini App | Open the bot [@aroma_type_test_bot](https://t.me/aroma_type_test_bot) in Telegram and launch the Mini App from the menu or inline button. |
| Production web | The frontend and API are served at [https://aromatypes.serveousercontent.com](https://aromatypes.serveousercontent.com). |
| Admin — Catalog Bot | A separate password-protected Telegram bot for adding, editing, toggling, and uploading photos of fragrances. Access requires the bot token and password (not yet transferred — see [Remaining Actions](#remaining-actions)). |
| Source code | Public GitHub repository: [https://github.com/BratishkaDurovaTg/SWP-AromaType](https://github.com/BratishkaDurovaTg/SWP-AromaType) (MIT License). |

### User Flow

1. Open the Mini App from Telegram.
2. Complete the 8-question psychotype questionnaire.
3. Receive up to 5 fragrance recommendations with explanations.
4. Open product cards to view fragrance details, notes, and accords.
5. Add fragrances to the cart, select volume (3 ml / 5 ml / 10 ml), and submit an order.
6. The order is sent to the configured Telegram contact.

### Admin Flow

1. Open the Catalog Bot (separate Telegram bot).
2. Authenticate with the configured password.
3. Use commands: `/add`, `/edit`, `/list`, `/view`, `/toggle`, `/photo` to manage fragrances.
4. Changes are reflected immediately in the public Mini App.

## Handover Scope — What Was Transferred, Delegated, or Retained

| Item | Status | Notes |
|---|---|---|
| Source repository | Public and accessible | The repository is public under MIT License. The customer can fork or clone at any time. |
| Source code and repository | To be transferred via private channel | The customer confirmed that source code (frontend, backend, admin components) is the primary deliverable. Repository access and a source code archive will be shared through a private channel. |
| Production VPS and domain | Retained by team — customer deprioritised | The customer explicitly confirmed that deployment activities are not a priority. The team retains server, DNS, and Caddy configuration. |
| Telegram Mini App bot (`@aroma_type_test_bot`) | Retained by team — to be transferred | The bot is registered under the team's Telegram account. BotFather settings, bot token, and Mini App URL will be transferred through a private channel. |
| Catalog admin bot | To be transferred | The catalog bot token (`CATALOG_BOT_TOKEN`) and password (`CATALOG_BOT_PASSWORD`) will be shared through a private channel. |
| Database | Retained by team | The PostgreSQL database runs on the team-managed VPS. No database access has been provided. |
| CI/CD pipelines | Retained by team | GitHub Actions workflows are configured under the team's GitHub organisation. |

## Required Configuration and Secrets-Handling

The production deployment uses environment variables. Create `/opt/aromatype/.env`
from the template file [`.env.production.example`](../.env.production.example) and
fill in the required values.

### Required Variables

| Variable | Purpose | Notes |
|---|---|---|
| `APP_ENV` | Runtime environment | Set to `production`. |
| `PORT` | Backend listen port | Default `8080`. |
| `DATABASE_URL` | PostgreSQL connection string | Must include user, password, host, port, and database name. |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | Must include the production domain. |
| `CATALOG_BOT_TOKEN` | Telegram bot token for the admin catalog bot | Obtain from BotFather. |
| `CATALOG_BOT_PASSWORD` | Password to authenticate admin commands | Set to a strong, unique value. |
| `POSTGRES_*` | PostgreSQL credentials | Must match the database configuration. |

### Secrets-Management Rules

- **Never commit `.env` files or production secrets to Git.** The `.gitignore` file
  already ignores `.env`.
- Use the `.env.production.example` as a template and replace every placeholder
  with real production values.
- The catalog bot is not started by the default `docker compose` command. Start it
  explicitly when needed:
  ```bash
  docker compose -f docker-compose.prod.yml --env-file .env --profile catalogbot up -d --build catalogbot
  ```

## Deployment and Setup Steps

The full deployment guide is maintained in [`docs/deployment.md`](deployment.md).
Below is a summary of the steps a new operator would follow:

1. **Provision a server** — Ubuntu 22.04 with Docker and the Compose plugin installed.
2. **Open firewall ports** — `22` (SSH), `80` (HTTP), `443` (HTTPS).
3. **Clone the repository** to the server (e.g., `/opt/aromatype`).
4. **Configure DNS** — Point the domain to the server's public IP.
5. **Create the environment file** — Copy `.env.production.example` to `/opt/aromatype/.env`
   and fill in all values.
6. **Start the services**:
   ```bash
   docker compose -f docker-compose.prod.yml --env-file .env up -d --build
   ```
7. **Verify the deployment**:
   ```bash
   curl -fsS https://aromatypes.serveousercontent.com/health
   ```
8. **Set the Mini App URL** in BotFather to `https://aromatypes.serveousercontent.com`.
9. **Start the catalog bot** (optional — only when catalog management is needed):
   ```bash
   docker compose -f docker-compose.prod.yml --env-file .env --profile catalogbot up -d --build catalogbot
   ```

### Recovery Steps

- **Check service status:** `docker compose -f docker-compose.prod.yml --env-file .env ps`
- **View logs:**
  ```bash
  docker compose -f docker-compose.prod.yml --env-file .env logs -f backend
  docker compose -f docker-compose.prod.yml --env-file .env logs -f catalogbot
  docker compose -f docker-compose.prod.yml --env-file .env logs -f caddy
  ```
- **Restart a service:** `docker compose -f docker-compose.prod.yml --env-file .env restart backend`
- **Rebuild and restart after code changes:** Re-run step 6 with `--build`.

## Operational Notes

- The back end serves a REST API documented in [`docs/api/openapi.yaml`](api/openapi.yaml).
  Swagger UI is available at `https://aromatypes.serveousercontent.com/docs`.
- The front end is a static single-page application served by Caddy. No build step
  is required on the server — Caddy serves the files directly from the repository.
- The recommendation engine is rule-based (psychotype scoring across four tags:
  drive, focus, aesthetic, power). No external AI service is required.
- The cart persists data in the browser's `localStorage`. Orders are sent to the
  configured Telegram contact — there is no persistent order management in the
  current version.
- Daily backups of the PostgreSQL database are recommended but not yet automated.
  The current setup does not include a backup mechanism.

## Troubleshooting and Support Guidance

| Symptom | Likely cause | Check |
|---|---|---|
| Mini App does not load | Mini App URL not configured in BotFather, or HTTPS not working | Verify the URL in BotFather; check Caddy logs. |
| API returns 5xx errors | Backend cannot connect to PostgreSQL, or environment misconfiguration | Check backend logs; verify `DATABASE_URL`. |
| Recommendations are empty | No active fragrances in the database, or psychotype tags not assigned | Use the catalog bot to add fragrances with psychotype scores. |
| Catalog bot not responding | `CATALOG_BOT_TOKEN` not set or bot not started | Verify the `.env` file; start the catalog bot with the `--profile catalogbot` flag. |
| Images not loading | `UPLOAD_DIR` misconfiguration or missing files | Check the upload volume mount; verify the `image_url` field. |

If the issue persists after checking the items above, the team can provide support
for the items that remain under team control (see [Remaining Actions](#remaining-actions)).

## Known Limitations

The following limitations and unfinished areas are known:

- **No automated database backups** — The current production deployment does not
  include a scheduled backup mechanism for the PostgreSQL database.
- **Administrator documentation pending** — A written guide for psychotype
  assignment and catalog bot usage is still in progress.
- **Minor UI refinements** — Several small UI improvements (input validation for
  phone/email, location display) were identified during earlier Sprints and
  remain in the backlog.

## Remaining Actions

| Action | Type | Blocks full transition? | Status |
|---|---|---|---|
| Transfer source code and repository access via private channel | Access/transfer | Yes — the customer needs the complete source code and admin components. | Pending — discussed, not yet executed |
| Transfer Telegram BotFather management and bot token | Account/ownership | Yes — the Mini App URL and bot token must be under customer control. | Pending — to be shared privately |
| Transfer catalog bot token and password | Access/credentials | Yes — the customer needs admin access to manage fragrances. | Pending — to be shared privately |
| Complete administrator documentation for psychotype assignment | Documentation | No — functional handover is accepted without it. | In progress |
| Set up automated database backups | Operational | No — customer deprioritised infrastructure. | Not started |
| Record public sanitized demo video for MVP v3 | Artifact | No — does not affect transition. | Not started |

### Blocker Classification

All remaining blockers are on the **team side** — the team has not yet completed
the transfer of access, credentials, and source code items. The product itself is
stable and ready for independent use. The customer confirmed during the Week 7
Sprint Review that deployment activities are not a priority and accepted the
current handover approach (source code transfer via private channel).

### Resolved During Sprint 5

The following issues from the Week 6 trial were fixed in Sprint 5:

| Issue | Resolution |
|---|---|
| No error pop-up in cart | Checkout validation now displays "Не все данные введены" when required fields are empty. |
| Duplicate product image in product card | Fixed — images no longer render twice. |
| Psychotype input in `/add` command | Redesigned — administrators now enter four comma-separated values totalling 100. |
| Cart items hidden behind checkout button | Cart scrolling improved so all items remain visible. |

## Documentation Entry Points

| Document | Purpose |
|---|---|
| [`README.md`](../README.md) | Main project entry point — what the product is, how to access it, links to all key docs. |
| [`docs/deployment.md`](deployment.md) | Production deployment guide — server setup, environment config, Docker Compose. |
| [`docs/testing.md`](testing.md) | Testing and QA status — CI gates, coverage, automated tests. |
| [`docs/quality-requirements.md`](quality-requirements.md) | Quality requirements (ISO/IEC 25010) with measurable scenarios. |
| [`docs/quality-requirement-tests.md`](quality-requirement-tests.md) | Automated quality requirement tests verifying each QR. |
| [`docs/user-acceptance-tests.md`](user-acceptance-tests.md) | UAT scenarios for customer-facing workflows. |
| [`docs/user-stories.md`](user-stories.md) | User-story index with stable IDs and traceability. |
| [`docs/roadmap.md`](roadmap.md) | Sprint-by-Sprint delivery plan. |
| [`docs/definition-of-done.md`](definition-of-done.md) | Team's minimum completion standard. |
| [`docs/api/openapi.yaml`](api/openapi.yaml) | OpenAPI specification for the public REST API. |
| [`docs/db-schema.md`](db-schema.md) | Database schema notes. |
| [`CHANGELOG.md`](../CHANGELOG.md) | Release history with user-visible changes per version. |

For normal customer use and operation, the most important documents are the
`README.md` (access and orientation), `docs/deployment.md` (setup and operations),
and the OpenAPI spec (API reference).

## Documentation Sufficiency Assessment

The current documentation set covers deployment, testing, API contracts, quality
requirements, user acceptance tests, and operational guidance. For the reached
handover level (**Ready for independent use**), the documentation is sufficient
for a technically capable operator to set up, deploy, and maintain the product,
provided that the access and ownership items in [Remaining Actions](#remaining-actions)
are completed.

Areas where documentation could be strengthened:
- Database backup and recovery procedures are not yet documented.
- Monitoring and alerting setup is not covered.
- Administrator documentation for psychotype assignment is still in progress.

These gaps do not block the current handover level but should be addressed before
full operational independence.
