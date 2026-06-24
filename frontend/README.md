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
In a deployed environment it uses the current origin unless another API URL is saved from the `API` button on the auth screen.

## MVP v1 Screens

- JWT login and registration.
- Onboarding screen.
- Guided questionnaire loaded from `GET /api/questions`.
- Recommendation result from `POST /api/recommendations`.
- Product card from `GET /api/fragrances/{id}`.
- Admin login and product creation form.
- Admin photo upload through `POST /api/admin/uploads/fragrance-photo`.

## Admin

For local Docker development, the default admin account is configured in `docker-compose.yml`:

```text
admin@aromatype.local
local-admin-password
```
