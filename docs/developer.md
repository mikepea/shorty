# Developer Guide

This guide covers setting up a development environment and contributing to Shorty.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Running Tests](#running-tests)
- [Code Style](#code-style)
- [Adding New Features](#adding-new-features)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)

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

## Development Setup

### Clone the Repository

```bash
git clone https://github.com/mikepea/shorty.git
cd shorty
```

### Install Development Tools

```bash
make tools
```

This installs:
- `swag` - Swagger documentation generator

### Start the Backend

```bash
go run ./cmd/shorty-server
```

The server starts on `http://localhost:8080`.

### Start the Frontend (Development Mode)

```bash
cd web
npm install
npm run dev
```

The frontend starts on `http://localhost:5173` with hot reloading. API requests are proxied to the backend.

### Default Login

- Email: `admin@shorty.local`
- Password: `changeme`

## Project Structure

```
shorty/
├── api/
│   └── swagger/           # Generated OpenAPI/Swagger docs
├── cmd/
│   └── shorty-server/     # Main application entry point
├── docs/                  # Documentation (you are here)
├── pkg/shorty/
│   ├── admin/             # Admin API handlers
│   ├── apikeys/           # API key authentication
│   ├── auth/              # User authentication (JWT)
│   ├── database/          # Database connection
│   ├── groups/            # Group management
│   ├── importexport/      # Bulk import/export
│   ├── links/             # Link management (core feature)
│   ├── models/            # GORM database models
│   ├── oidc/              # OIDC/SSO integration
│   ├── redirect/          # URL redirect handler
│   ├── scim/              # SCIM 2.0 provisioning
│   └── tags/              # Tag management
├── tests/
│   └── integration/       # Integration tests
├── web/                   # React frontend
│   ├── src/
│   │   ├── api/           # API client
│   │   ├── components/    # Reusable UI components
│   │   ├── context/       # React context (auth state)
│   │   └── pages/         # Page components
│   └── package.json
├── Makefile               # Build and development tasks
├── go.mod                 # Go dependencies
└── go.sum
```

### Key Packages

| Package | Purpose |
|---------|---------|
| `pkg/shorty/models` | Database models and migrations |
| `pkg/shorty/auth` | JWT authentication and password hashing |
| `pkg/shorty/links` | Core link shortening logic |
| `pkg/shorty/scim` | SCIM 2.0 user/group provisioning |
| `pkg/shorty/oidc` | OpenID Connect SSO |

## Running Tests

### All Tests

```bash
# Using Make
make test

# Or directly
go test -v ./pkg/...
go test -v ./tests/integration/...
```

### Specific Package

```bash
go test -v ./pkg/shorty/links/...
```

### With Coverage

```bash
go test -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out
```

### Integration Tests with Keycloak

```bash
INTEGRATION_TEST_KEYCLOAK=1 go test -v ./tests/integration/...
```

This requires Docker to run a Keycloak container.

## Code Style

### Go

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing

```bash
gofmt -w .
go vet ./...
```

### Frontend

- Use TypeScript
- Follow React best practices
- Use functional components with hooks

```bash
cd web
npm run lint
```

## Adding New Features

### Adding a New API Endpoint

1. **Create or update the handler** in the appropriate package under `pkg/shorty/`:

```go
// pkg/shorty/example/handlers.go
package example

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{db: db}
}

// MyEndpoint does something
// @Summary Short description
// @Description Longer description
// @Tags example
// @Produce json
// @Success 200 {object} MyResponse
// @Security BearerAuth
// @Router /example [get]
func (h *Handler) MyEndpoint(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "hello"})
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
    rg.GET("/example", h.MyEndpoint)
}
```

2. **Register routes** in `cmd/shorty-server/main.go`:

```go
exampleHandler := example.NewHandler(database.GetDB())
exampleHandler.RegisterRoutes(api.Group("", combinedAuth))
```

3. **Regenerate Swagger docs**:

```bash
make generate
```

4. **Write tests** in `pkg/shorty/example/handlers_test.go`

5. **Update frontend** if needed in `web/src/`

### Adding a New Database Model

1. **Define the model** in `pkg/shorty/models/`:

```go
// pkg/shorty/models/example.go
package models

import "gorm.io/gorm"

type Example struct {
    gorm.Model
    Name        string `gorm:"not null"`
    Description string
    UserID      uint   `gorm:"not null"`
    User        User   `gorm:"foreignKey:UserID"`
}
```

2. **Add to migrations** in `pkg/shorty/models/models.go`:

```go
func AllModels() []interface{} {
    return []interface{}{
        // ... existing models
        &Example{},
    }
}
```

3. **The migration runs automatically** on server startup.

## API Documentation

### Swagger Annotations

Add Swagger annotations to handler functions:

```go
// MyHandler does something
// @Summary Short description
// @Description Longer description of what this endpoint does
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "Resource ID"
// @Param request body MyRequest true "Request body"
// @Success 200 {object} MyResponse
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /resource/{id} [put]
func (h *Handler) MyHandler(c *gin.Context) {
    // ...
}
```

### Regenerating Documentation

After adding or modifying Swagger annotations:

```bash
make generate
```

This updates `api/swagger/swagger.json` and `api/swagger/swagger.yaml`.

### Viewing Documentation

Start the server and visit: `http://localhost:8080/swagger/index.html`

## Contributing

### Git Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes with clear commits

3. Ensure tests pass:
   ```bash
   make test
   ```

4. Ensure generated files are up to date:
   ```bash
   make check-generate
   ```

5. Push and create a PR:
   ```bash
   git push -u origin feature/my-feature
   gh pr create
   ```

### Commit Message Format

```
Short description of change

Optional longer description if needed.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Your Name <your@email.com>
```

### Pull Request Guidelines

- PRs require passing CI checks before merge
- Include tests for new functionality
- Update documentation if needed
- Keep PRs focused on a single feature or fix

### CI Checks

The CI pipeline runs:

1. **Go tests** - Unit and integration tests
2. **Generate check** - Ensures Swagger docs are up to date
3. **Frontend build** - Ensures React app compiles

All checks must pass before merging.

## Debugging

### Server Logs

The server logs to stdout. Useful log messages include:
- Database connection status
- Migration results
- Request/response logging (in debug mode)

### Database Queries

Enable GORM debug mode to see SQL queries:

```go
db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

### Frontend Debugging

Use React Developer Tools browser extension and the Network tab to debug API calls.
