# Static View — AromaType

The static view describes the system's component structure and module boundaries
at the current delivered state (MVP v3).

## Component Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                   Telegram Mini App (Frontend)                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────┐ │
│  │  Home    │  │  Quiz    │  │ Results  │  │    Cart /     │ │
│  │ Screen   │  │ Screen   │  │ Screen   │  │   Checkout    │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────┘ │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              API Client (fetch)                          │ │
│  └──────────────────────────────────────────────────────────┘ │
└───────────────────────┬──────────────────────────────────────┘
                        │ HTTP (JSON)
                        ▼
┌──────────────────────────────────────────────────────────────┐
│                     Go API Server (:8080)                     │
│  ┌──────────┐  ┌──────────┐  ┌────────────────────────────┐ │
│  │  CORS    │  │ Logging  │  │       http.Router          │ │
│  │Middleware│  │Middleware│  │  ┌──────────────────────┐  │ │
│  └──────────┘  └──────────┘  │  │ /health              │  │ │
│                               │  │ /api/questions       │  │ │
│                               │  │ /api/recommendations │  │ │
│                               │  │ /api/fragrances/{id} │  │ │
│                               │  │ /docs, /openapi.yaml │  │ │
│                               │  │ /uploads/            │  │ │
│                               │  └──────────────────────┘  │ │
│                               └────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              questionnaire.Service                        │ │
│  │  ┌─────────────────┐  ┌────────────────────────────────┐ │ │
│  │  │  Recommendation │  │  Fragrance CRUD                │ │ │
│  │  │  Engine         │  │  (Create/Read)                 │ │ │
│  │  └─────────────────┘  └────────────────────────────────┘ │ │
│  └──────────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              questionnaire.Repository                     │ │
│  │  (pgxpool connection to PostgreSQL)                       │ │
│  └──────────────────────────────────────────────────────────┘ │
└───────────────────────────────┬──────────────────────────────┘
                                │ SQL
                                ▼
┌──────────────────────────────────────────────────────────────┐
│                      PostgreSQL 16                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐  ┌──────────┐ │
│  │questions │  │answer_   │  │fragrances     │  │  tags    │ │
│  │          │  │options   │  │(JSONB fields) │  │          │ │
│  └──────────┘  └──────────┘  └──────────────┘  └──────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌────────────────────────────┐ │
│  │fragrance_│  │answer_   │  │  (seed data: 5 fragrances, │ │
│  │tags      │  │option_   │  │   8 questions, 32 options,  │ │
│  │          │  │tags      │  │   24 tags)                  │ │
│  └──────────┘  └──────────┘  └────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│                   Catalog Bot (Separate Binary)               │
│  ┌──────────┐  ┌──────────┐  ┌────────────────────────────┐ │
│  │Telegram  │  │ Session  │  │  catalogbot Handlers       │ │
│  │Client    │  │ Manager  │  │  ┌──────────────────────┐  │ │
│  │(longpoll)│  │(per-chat)│  │  │ /start, /list, /add  │  │ │
│  └──────────┘  └──────────┘  │  │ /edit, /view, /photo │  │ │
│                               │  │ /toggle, /set, /help │  │ │
│                               │  └──────────────────────┘  │ │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │              catalog.Repository                           │ │
│  │  (List, Get, Upsert fragrances)                           │ │
│  └──────────────────────────────────────────────────────────┘ │
└───────────────────────────────┬──────────────────────────────┘
                                │ SQL
                                ▼
                        PostgreSQL 16
```

## Module Descriptions

### Frontend (`frontend/`)

| Module | Responsibility |
|---|---|
| [`index.html`](../../frontend/index.html) | SPA shell — loads Telegram WebApp SDK, app container, toast container |
| [`app.js`](../../frontend/app.js) | Event-driven SPA with hash-based routing, 11 screens, API client, cart and order persistence |
| [`styles.css`](../../frontend/styles.css) | Pastel-colour design system, premium styling |

### Backend (`backend/internal/`)

| Package | Responsibility | Key Types |
|---|---|---|
| `http` | HTTP server, routing, CORS, logging middleware | `Router`, `withCORS` |
| `questionnaire` | Core business logic — recommendation engine, domain types | `Service`, `Repository`, `Fragrance`, `Profile`, `RecommendationResponse` |
| `catalog` | Fragrance CRUD repository for catalog management | `Repository` |
| `catalogbot` | Telegram admin bot — session management, 16-step add wizard, input parsing | `Bot`, `session`, `telegramClient` |
| `config` | Environment-based configuration | `Config` |
| `database` | PostgreSQL connection, migration, seed data | `Connect`, `Migrate`, `seedMVPData` |

### Entry Points

| Binary | File | Responsibility |
|---|---|---|
| `aromatype-api` | `cmd/api/main.go` | HTTP server — config, DB connect, migrate, HTTP listener |
| `aromatype-catalog-bot` | `cmd/catalogbot/main.go` | Catalog bot — config, DB connect, migrate, long-poll loop |

## Module Dependencies

```
cmd/api/main.go
  └── config
  └── database
  └── questionnaire.Service ──┬── questionnaire.Repository ──> PostgreSQL
                               └── http.Router ──> HTTP

cmd/catalogbot/main.go
  └── config
  └── database
  └── catalogbot.Bot ──┬── catalog.Repository ──> PostgreSQL
                        └── telegramClient ──> Telegram API
```

The questionnaire `Service` depends on a `fragranceRepository` interface, not
on the concrete `Repository` — enabling unit testing with fake implementations.
