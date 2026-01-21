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
