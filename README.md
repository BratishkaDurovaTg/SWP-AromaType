# SWP-AromaType

AromaType is a Telegram Mini App for personalized fragrance discovery through
style, situations, and feelings rather than perfume terminology.

## Repository Structure

```text
backend/   Go API, PostgreSQL migrations, recommendation logic
frontend/  Web client and admin UI
docs/      API contract, database schema notes, product docs
```

## Local Development

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

## Team Workflow

- `main` stores stable versions.
- Work in feature branches, for example `feature/backend-auth`.
- Keep API changes documented in `docs/api/openapi.yaml`.
- Frontend can use mock data until backend endpoints are ready.
