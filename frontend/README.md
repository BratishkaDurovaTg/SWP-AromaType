# Aroma Type Frontend

Static Telegram Mini App frontend for Aroma Type MVP v1.

## Run Locally

Start the backend first:

```bash
docker compose up --build
```

Start the frontend from this directory:

```bash
python3 -m http.server 5173
```

Open:

```text
http://localhost:5173
```

The local frontend uses `http://localhost:8080` as API base by default.
In a deployed environment it uses the current origin.

## MVP v1 Screens

- Onboarding screen.
- Guided questionnaire loaded from `GET /api/questions`.
- Recommendation result from `POST /api/recommendations`.
- Product card from `GET /api/fragrances/{id}`.
- Telegram contact handoff for ordering a sample set or a fragrance.
