# Deployment View — AromaType

The deployment view describes the infrastructure, network topology, and hosting
arrangements for both local development and production environments.

## Development Deployment

```
┌─────────────────────────────────────────────────┐
│                Developer Machine                 │
│                                                   │
│  ┌──────────────────────┐  ┌──────────────────┐  │
│  │   Frontend Dev Server │  │  Browser / Phone  │  │
│  │   (python3 http.srv)  │  │  (Telegram Mini   │  │
│  │   Port 5173           │  │   App client)     │  │
│  └──────────────────────┘  └──────────────────┘  │
│           │                        │              │
│           └────────┬───────────────┘              │
│                    │ HTTP                          │
│                    ▼                               │
│  ┌─────────────────────────────────────────────┐  │
│  │           Docker Compose                     │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  │  │
│  │  │ backend  │  │catalogbot│  │ postgres │  │  │
│  │  │ :8080    │  │(profile  │  │ :5432    │  │  │
│  │  │          │  │ bot)     │  │          │  │  │
│  │  └──────────┘  └──────────┘  └──────────┘  │  │
│  └─────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

**Start:**
```bash
docker compose up --build
cd frontend && python3 -m http.server 5173
```

## Production Deployment

```
                           Internet
                              │
                              ▼
              ┌──────────────────────────────┐
              │   Ubuntu 22.04 VPS            │
              │                               │
              │  ┌────────────────────────┐   │
              │  │      Caddy 2            │   │
              │  │  Port 80 → 443 (HTTPS)  │   │
              │  │  Let's Encrypt (auto)   │   │
              │  │                         │   │
              │  │  ┌── /api/*  ─────┐     │   │
              │  │  │  /health      │     │   │
              │  │  │  /docs        │──┐  │   │
              │  │  │  /openapi.yaml│  │  │   │
              │  │  │  /uploads/*   │  │  │   │
              │  │  └───────────────┘  │  │   │
              │  │  ┌──────────────┐   │  │   │
              │  │  │  / (SPA)     │◄──┘  │   │
              │  │  │  index.html  │      │   │
              │  │  └──────────────┘      │   │
              │  └────────────────────────┘   │
              │           │                   │
              │           ▼                   │
              │  ┌──────────────────┐         │
              │  │   Go API Server  │         │
              │  │   Port 8080      │         │
              │  └──────────────────┘         │
              │           │                   │
              │           ▼                   │
              │  ┌──────────────────┐         │
              │  │   PostgreSQL 16  │         │
              │  │   Port 5432      │         │
              │  └──────────────────┘         │
              │                               │
              │  ┌──────────────────────────┐ │
              │  │   Catalog Bot (optional) │ │
              │  │   Long-poll to Telegram  │ │
              │  │   Profile: catalogbot    │ │
              │  └──────────────────────────┘ │
              └───────────────────────────────┘
```

**Production URL:** [https://aromatypes.serveousercontent.com](https://aromatypes.serveousercontent.com)

**Deploy:**
```bash
docker compose -f docker-compose.prod.yml --env-file .env up -d --build
```

## Container Configuration

### Docker Compose Services

| Service | Image | Port(s) | Dependencies | Volume(s) |
|---|---|---|---|---|
| `backend` | Build from `backend/Dockerfile` | 8080 | postgres (healthy) | `uploads_data` |
| `catalogbot` | Build from `backend/Dockerfile` | — | postgres | — |
| `postgres` | `postgres:16-alpine` | 5432 | — | `postgres_data` |
| `caddy` | `caddy:2-alpine` | 80, 443 | backend | `caddy_data`, `caddy_config`, frontend files |

### Environment Variables

| Variable | Source | Used By |
|---|---|---|
| `APP_ENV` | `.env` | backend |
| `DATABASE_URL` | `.env` | backend, catalogbot |
| `PORT` | `.env` | backend |
| `CORS_ALLOWED_ORIGINS` | `.env` | backend |
| `CATALOG_BOT_TOKEN` | `.env` | catalogbot |
| `CATALOG_BOT_PASSWORD` | `.env` | catalogbot |
| `POSTGRES_*` | `.env` | postgres |

## Network Flows

```
External → Caddy (443) → backend (8080): API requests
External → Caddy (443) → static files:  Frontend SPA
Catalog Bot → Telegram API:              Long-poll updates (outbound only)
Backend → PostgreSQL (5432):             Database queries
```

## CI/CD Pipeline

```
GitHub Push/PR
  │
  ├── backend.yml:  gofmt, go vet, go test -race, coverage
  ├── frontend-tests.yml: npm ci, npm test
  ├── qa.yml: govulncheck
  ├── docker-build.yml: docker build backend
  └── lychee.yml: markdown link check
  │
  └── deploy-backend.yml (manual):
        SSH → VPS → git pull → docker compose up -d --build
```

## Security Boundaries

- The frontend is served over HTTPS (Caddy auto TLS).
- The API is stateless — no authentication, no JWT, no session tokens.
- The catalog bot requires a password for admin operations.
- The PostgreSQL database is not exposed externally — only accessible via
  the Docker Compose internal network.
- Uploaded images are served through the API at `/uploads/`.
