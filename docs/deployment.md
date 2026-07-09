# Customer Deployment Guide

This guide describes how to deploy AromaType on a customer-owned VPS with a
customer-owned domain.

The production stack runs as one Docker Compose product:

- Caddy reverse proxy and HTTPS termination
- Static Telegram Mini App frontend
- Go backend API
- PostgreSQL database
- Optional separate Telegram catalog admin bot

## 1. Requirements

VPS:

- Ubuntu 22.04 LTS or newer
- 1 vCPU / 2 GB RAM minimum
- Public IPv4 address
- Open inbound ports `22`, `80`, and `443`

Domain:

- A domain controlled by the customer
- DNS access for creating `A` records

Local access:

- SSH access to the VPS as `root` or a sudo-capable user
- GitHub access to clone this repository

Telegram:

- Public customer-facing bot for launching the Mini App
- Separate catalog admin bot for product management

## 2. DNS Setup

Create DNS records at the customer's registrar or DNS provider:

```text
Type  Name  Value
A     @     <VPS_PUBLIC_IP>
A     www   <VPS_PUBLIC_IP>
```

Wait until DNS resolves:

```bash
dig +short customer-domain.example
dig +short www.customer-domain.example
```

Both commands should return the VPS public IP.

## 3. Install Server Dependencies

Connect to the VPS:

```bash
ssh root@<VPS_PUBLIC_IP>
```

Install Docker, Docker Compose plugin, Git, and firewall rules:

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

Verify Docker:

```bash
docker --version
docker compose version
```

## 4. Clone The Project

Use `/opt/aromatype` as the production directory:

```bash
mkdir -p /opt
cd /opt
git clone https://github.com/BratishkaDurovaTg/SWP-AromaType.git aromatype
cd /opt/aromatype
```

Checkout the branch or tag approved for production:

```bash
git checkout main
git pull origin main
```

## 5. Configure Environment

Create the production environment file:

```bash
cp .env.production.example .env
nano .env
```

Required values:

```text
APP_DOMAIN=customer-domain.example
APP_WWW_DOMAIN=www.customer-domain.example

POSTGRES_DB=aromatype
POSTGRES_USER=aromatype
POSTGRES_PASSWORD=<strong-random-password>

CORS_ALLOWED_ORIGINS=https://customer-domain.example,https://www.customer-domain.example

CATALOG_BOT_TOKEN=<telegram-catalog-admin-bot-token>
CATALOG_BOT_PASSWORD=<strong-admin-password>
```

Rules:

- `APP_DOMAIN` and `APP_WWW_DOMAIN` must match DNS records.
- `CORS_ALLOWED_ORIGINS` must include every public HTTPS origin that will call
  the backend.
- Do not commit `.env`.
- Use a separate Telegram bot token for catalog administration.

## 6. Start Production Services

Start the public product without the catalog bot:

```bash
docker compose -f docker-compose.prod.yml --env-file .env up -d --build
```

Start the catalog admin bot as well:

```bash
docker compose -f docker-compose.prod.yml --env-file .env --profile catalogbot up -d --build
```

Check containers:

```bash
docker compose -f docker-compose.prod.yml --env-file .env ps
```

## 7. Verify Deployment

Health check:

```bash
curl -fsS https://customer-domain.example/health
```

Questionnaire API:

```bash
curl -fsS https://customer-domain.example/api/questions
```

Swagger UI:

```text
https://customer-domain.example/docs
```

Frontend:

```text
https://customer-domain.example
```

Logs:

```bash
docker compose -f docker-compose.prod.yml --env-file .env logs -f caddy
docker compose -f docker-compose.prod.yml --env-file .env logs -f backend
docker compose -f docker-compose.prod.yml --env-file .env logs -f catalogbot
```

## 8. Configure Telegram Mini App

In BotFather, configure the customer-facing bot:

```text
Mini App URL: https://customer-domain.example
```

If the bot uses a menu button, configure the same URL as the Web App button URL.

Important:

- Telegram Mini Apps require HTTPS.
- Do not use the catalog admin bot as the customer-facing bot.
- The catalog admin bot is only for product management.

## 9. Catalog Admin Bot Smoke Test

Open the catalog admin bot in Telegram:

```text
/start
```

Enter `CATALOG_BOT_PASSWORD`.

Expected result:

- The bot opens the catalog admin menu.
- `Каталог` lists products.
- `Добавить товар` starts the product creation flow.

## 10. Updating Production

Deploy a new approved release:

```bash
cd /opt/aromatype
git fetch origin
git checkout main
git pull origin main
docker compose -f docker-compose.prod.yml --env-file .env --profile catalogbot up -d --build
docker compose -f docker-compose.prod.yml --env-file .env ps
```

Run a smoke test:

```bash
curl -fsS https://customer-domain.example/health
curl -fsS https://customer-domain.example/api/questions
```

## 11. Rollback

Find the previous stable commit:

```bash
git log --oneline -10
```

Checkout it and rebuild:

```bash
git checkout <previous-stable-commit>
docker compose -f docker-compose.prod.yml --env-file .env --profile catalogbot up -d --build
```

After the rollback, verify:

```bash
curl -fsS https://customer-domain.example/health
docker compose -f docker-compose.prod.yml --env-file .env ps
```

## 12. Backup Notes

PostgreSQL data is stored in the Docker volume `aromatype_postgres_data`.
Product images are stored in `aromatype_uploads_data`.

Before destructive maintenance, create backups:

```bash
docker compose -f docker-compose.prod.yml --env-file .env exec postgres pg_dump -U "$POSTGRES_USER" "$POSTGRES_DB" > aromatype-backup.sql
docker run --rm -v aromatype_uploads_data:/data -v "$PWD":/backup alpine tar czf /backup/aromatype-uploads.tar.gz /data
```

Store backups outside the VPS.
