# Development Environment Setup

This guide covers setting up your local development environment for Shorty.

## Prerequisites

- Go 1.25+
- Node.js 20+
- Git
- Make (optional, but recommended)

### Installing Go

```bash
# macOS with Homebrew
brew install go

# Or use asdf
asdf plugin add golang
asdf install golang 1.25.6
asdf global golang 1.25.6
```

### Installing Node.js

```bash
# macOS with Homebrew
brew install node@20

# Or use asdf
asdf plugin add nodejs
asdf install nodejs 20.10.0
asdf global nodejs 20.10.0
```

## Clone the Repository

```bash
git clone https://github.com/mikepea/shorty.git
cd shorty
```

## Install Development Tools

```bash
make tools
```

This installs:
- `swag` - Swagger documentation generator

## Start the Backend

```bash
go run ./cmd/shorty-server
```

The server starts on `http://localhost:8080`.

## Start the Frontend (Development Mode)

```bash
cd web
npm install
npm run dev
```

The frontend starts on `http://localhost:3000` with hot reloading. API requests are proxied to the backend.

## Default Login

When starting with a fresh database:

- Email: `admin@shorty.local`
- Password: `changeme`

## Docker Setup (Alternative)

You can run Shorty with Docker instead of installing Go and Node.js locally.

### Prerequisites

- Docker
- Docker Compose

### Quick Start

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f shorty

# Stop
docker-compose down
```

The app runs at `http://localhost:8080` with SQLite storage persisted in a Docker volume.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SHORTY_DB_PATH` | `/data/shorty.db` | SQLite database path |
| `JWT_SECRET` | (set in compose) | Secret for JWT signing |
| `SHORTY_BASE_URL` | `http://localhost:8080` | Base URL for generated links |
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `release` | Gin framework mode |

### Rebuilding After Code Changes

```bash
docker-compose build
docker-compose up -d
```

## Project Structure Overview

```
shorty/
├── api/
│   └── swagger/           # Generated OpenAPI/Swagger docs
├── cmd/
│   └── shorty-server/     # Main application entry point
├── docs/                  # Documentation
├── pkg/shorty/            # Backend Go packages
├── tests/
│   └── integration/       # Integration tests
├── web/                   # React frontend
├── Makefile               # Build and development tasks
├── go.mod                 # Go dependencies
└── go.sum
```

## Next Steps

- [Backend Development](backend.md) - Working with Go code
- [Frontend Development](frontend.md) - Working with React code
- [Debugging & Diagnostics](diagnosis.md) - Troubleshooting issues
