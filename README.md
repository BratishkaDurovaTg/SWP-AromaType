# SWP-AromaType

AromaType is a Telegram Mini App for personalized fragrance discovery through
style, situations, and feelings rather than perfume terminology.

## Quick Access

- **Telegram Mini App:** [@aroma_type_test_bot](https://t.me/aroma_type_test_bot)
- **Production deployment:** [https://aromatypes.serveousercontent.com](https://aromatypes.serveousercontent.com)
- **Hosted documentation:** [GitHub Pages](https://github.com/BratishkaDurovaTg/SWP-AromaType/deployments/github-pages)
- **Source code:** [GitHub repository](https://github.com/BratishkaDurovaTg/SWP-AromaType)
- **License:** MIT

## Repository Structure

```text
backend/   Go API, PostgreSQL migrations, recommendation logic
frontend/  Telegram Mini App web client
docs/      API contract, database schema notes, product docs
```

## Documentation

- [Contributing guide](CONTRIBUTING.md)
- [Agent guidance](AGENTS.md)
- [Customer handover](docs/customer-handover.md)
- [User stories](docs/user-stories.md)
- [Definition of Done](docs/definition-of-done.md)
- [Roadmap](docs/roadmap.md)
- [Testing and QA status](docs/testing.md)
- [Quality requirements](docs/quality-requirements.md)
- [Quality requirement tests](docs/quality-requirement-tests.md)
- [User acceptance tests](docs/user-acceptance-tests.md)
- [Deployment guide](docs/deployment.md)
- [API contract](docs/api/openapi.yaml)
- [Database schema](docs/db-schema.md)
- [Catalog bot](docs/catalog-bot.md)

## Local Development

### Full stack (frontend + backend)

Requirements:

- Docker Desktop
- Python 3

Run the backend and PostgreSQL from the repository root:

```bash
docker compose up --build
```

In a second terminal, run the frontend:

```bash
cd frontend
python3 -m http.server 5173
```

Open the app:

```text
http://localhost:5173
```

Useful local links:

- Frontend: `http://localhost:5173`
- Backend health: `http://localhost:8080/health`
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI spec: `http://localhost:8080/openapi.yaml`

### Backend only

Start the backend stub:

```bash
cd backend
go run ./cmd/api
```

Open:

- API health check: `/health`
- Swagger UI: `/docs`
- OpenAPI spec: `/openapi.yaml`

Start with Docker Compose:

```bash
docker compose up --build
```

PostgreSQL connection for local backend development:

```text
host: localhost
port: 5432
database: aromatype
user: aromatype
password: aromatype
```

The public MVP API does not require registration or JWT. Users complete the
questionnaire anonymously, receive up to 5 recommended fragrances, and order
through the configured Telegram contact. Catalog management is handled through a
separate password-protected Telegram bot workflow, not through this web client.

## Team Workflow

- `main` and `dev` are protected branches — no direct pushes.
- Create branches from the relevant issue: `<issue-number>-short-description`.
- Submit changes through a Pull Request and obtain approval from another team member.
- Use merge commits (squash and rebase are disabled).
- Keep API changes documented in `docs/api/openapi.yaml`.
- See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full contribution workflow.
