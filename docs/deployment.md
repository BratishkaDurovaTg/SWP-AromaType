# Deployment

Production deployment uses Docker Compose with Caddy, the Go backend, and PostgreSQL.

## Domain

The current production domain is:

- `https://aroma-type.shop`
- `https://www.aroma-type.shop`

Both records must point to the VPS public IP before Caddy can issue HTTPS certificates.

## Server Setup

Install Docker and the Compose plugin on Ubuntu 22.04:

```bash
apt-get update
apt-get install -y ca-certificates curl git ufw
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" > /etc/apt/sources.list.d/docker.list
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
ufw allow OpenSSH
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable
```

## Environment

Create `/opt/aromatype/.env` from [.env.production.example](../.env.production.example) and replace every placeholder with production values.

Do not commit production secrets to Git.

## Deploy

From `/opt/aromatype`:

```bash
docker compose -f docker-compose.prod.yml --env-file .env up -d --build
```

Useful checks:

```bash
docker compose -f docker-compose.prod.yml --env-file .env ps
docker compose -f docker-compose.prod.yml --env-file .env logs -f backend
docker compose -f docker-compose.prod.yml --env-file .env logs -f caddy
curl -fsS https://aroma-type.shop/health
```

## Telegram Mini App

After HTTPS is working, set the Mini App URL in BotFather:

```text
https://aroma-type.shop
```

If using a Telegram menu button, configure the same URL as the Web App button URL.
