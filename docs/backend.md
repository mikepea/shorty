# Backend Development

This guide covers developing the Go backend for Shorty.

## Project Structure

```
pkg/shorty/
├── admin/             # Admin API handlers
├── apikeys/           # API key authentication
├── auth/              # User authentication (JWT)
├── database/          # Database connection
├── groups/            # Group management
├── importexport/      # Bulk import/export
├── links/             # Link management (core feature)
├── models/            # GORM database models
├── oidc/              # OIDC/SSO integration
├── redirect/          # URL redirect handler
├── scim/              # SCIM 2.0 provisioning
└── tags/              # Tag management
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

### All Backend Tests

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

- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` before committing

```bash
gofmt -w .
go vet ./...
```

## Adding a New API Endpoint

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

## Adding a New Database Model

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

## CI Checks

The CI pipeline runs these backend checks:

1. **Go tests** - Unit and integration tests
2. **Generate check** - Ensures Swagger docs are up to date
3. **Build** - Ensures the code compiles

Ensure generated files are up to date before pushing:

```bash
make check-generate
```
