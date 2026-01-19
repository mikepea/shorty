# Deployment Guide

This guide covers deploying Shorty in production environments.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Reverse Proxy Configuration](#reverse-proxy-configuration)
- [Production Checklist](#production-checklist)

## Prerequisites

- Go 1.25+ (for building from source)
- Node.js 20+ (for building frontend)
- PostgreSQL 14+ (recommended for production) or SQLite
- A reverse proxy (nginx, Caddy, or similar)
- TLS certificate (Let's Encrypt recommended)

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `DATABASE_URL` | Database connection string | `shorty.db` | Yes |
| `JWT_SECRET` | Secret for signing JWT tokens | Auto-generated | **Yes** |
| `SHORTY_BASE_URL` | Public URL (for OIDC callbacks, SCIM) | `http://localhost:8080` | Yes |

### JWT_SECRET

Generate a secure random secret:

```bash
openssl rand -base64 32
```

**Important:** Set this explicitly in production. If not set, a random secret is generated on each restart, invalidating all existing tokens.

## Database Setup

### PostgreSQL (Recommended)

1. Create a database and user:

```sql
CREATE USER shorty WITH PASSWORD 'your-secure-password';
CREATE DATABASE shorty OWNER shorty;
```

2. Set the connection string:

```bash
export DATABASE_URL="postgres://shorty:your-secure-password@localhost:5432/shorty?sslmode=require"
```

### SQLite

SQLite works for small deployments but is not recommended for production:

```bash
export DATABASE_URL="/var/lib/shorty/shorty.db"
```

Ensure the directory exists and is writable by the application user.

## Docker Deployment

### Using Docker Compose

Create a `docker-compose.yml`:

```yaml
version: '3.8'

services:
  shorty:
    image: ghcr.io/mikepea/shorty:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://shorty:changeme@db:5432/shorty?sslmode=disable
      - JWT_SECRET=your-secure-jwt-secret-here
      - SHORTY_BASE_URL=https://go.example.com
    depends_on:
      - db
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=shorty
      - POSTGRES_PASSWORD=changeme
      - POSTGRES_DB=shorty
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  postgres_data:
```

Start the services:

```bash
docker-compose up -d
```

### Building the Docker Image

```dockerfile
# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o shorty-server ./cmd/shorty-server

# Frontend build
FROM node:20-alpine AS frontend
WORKDIR /app
COPY web/package*.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /app/shorty-server .
COPY --from=frontend /app/dist ./web/dist

EXPOSE 8080
CMD ["./shorty-server"]
```

Build and run:

```bash
docker build -t shorty .
docker run -d -p 8080:8080 \
  -e DATABASE_URL="postgres://..." \
  -e JWT_SECRET="..." \
  -e SHORTY_BASE_URL="https://go.example.com" \
  shorty
```

## Manual Deployment

### Building from Source

```bash
# Clone the repository
git clone https://github.com/mikepea/shorty.git
cd shorty

# Build the server
go build -o shorty-server ./cmd/shorty-server

# Build the frontend
cd web
npm ci
npm run build
cd ..
```

### Systemd Service

Create `/etc/systemd/system/shorty.service`:

```ini
[Unit]
Description=Shorty URL Shortener
After=network.target postgresql.service

[Service]
Type=simple
User=shorty
Group=shorty
WorkingDirectory=/opt/shorty
ExecStart=/opt/shorty/shorty-server
Restart=always
RestartSec=5

Environment=PORT=8080
Environment=DATABASE_URL=postgres://shorty:password@localhost:5432/shorty?sslmode=disable
Environment=JWT_SECRET=your-secure-jwt-secret
Environment=SHORTY_BASE_URL=https://go.example.com

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable shorty
sudo systemctl start shorty
```

## Reverse Proxy Configuration

### Nginx

```nginx
server {
    listen 80;
    server_name go.example.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name go.example.com;

    ssl_certificate /etc/letsencrypt/live/go.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/go.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Caddy

```caddyfile
go.example.com {
    reverse_proxy localhost:8080
}
```

Caddy automatically provisions TLS certificates via Let's Encrypt.

## Production Checklist

Before going live, ensure:

- [ ] **JWT_SECRET** is set to a secure, random value
- [ ] **DATABASE_URL** points to PostgreSQL (not SQLite)
- [ ] **SHORTY_BASE_URL** matches your public URL
- [ ] TLS is configured (HTTPS only)
- [ ] Database backups are configured
- [ ] Log aggregation is set up
- [ ] Monitoring/health checks are configured (`/health` endpoint)
- [ ] Default admin password has been changed
- [ ] Firewall rules allow only necessary ports (443, 80 for redirect)

### Health Check

Shorty exposes a health endpoint:

```bash
curl https://go.example.com/health
# {"status":"ok"}
```

Use this for load balancer health checks and monitoring.

### Backup Strategy

For PostgreSQL:

```bash
# Daily backup
pg_dump -U shorty shorty > /backups/shorty-$(date +%Y%m%d).sql

# Restore
psql -U shorty shorty < /backups/shorty-20240115.sql
```

### Logging

Shorty logs to stdout. In production, configure your container runtime or systemd to capture and rotate logs:

```bash
# View logs with systemd
journalctl -u shorty -f

# View logs with Docker
docker logs -f shorty
```
