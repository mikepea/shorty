# Developer Guide

Welcome to the Shorty developer documentation. This guide will help you set up your development environment and contribute to the project.

## Documentation

| Guide | Description |
|-------|-------------|
| [Setup](setup.md) | Development environment setup, prerequisites, and first run |
| [Backend](backend.md) | Go backend development, API endpoints, and database models |
| [Frontend](frontend.md) | React frontend development, testing, and components |
| [Debugging](diagnosis.md) | Troubleshooting, logging, and diagnostics |

## Quick Start

```bash
# Clone and setup
git clone https://github.com/mikepea/shorty.git
cd shorty
make tools

# Start backend (terminal 1)
go run ./cmd/shorty-server

# Start frontend (terminal 2)
cd web && npm install && npm run dev
```

Default login: `admin@shorty.local` / `changeme`

## Contributing

### Git Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Make your changes with clear commits

3. Ensure tests pass:
   ```bash
   # Backend
   go test -v ./pkg/...

   # Frontend
   cd web && npm test
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
```

### Pull Request Guidelines

- PRs require passing CI checks before merge
- Include tests for new functionality
- Update documentation if needed
- Keep PRs focused on a single feature or fix

### CI Checks

The CI pipeline runs:

1. **Backend tests** - Go unit and integration tests
2. **Generate check** - Ensures Swagger docs are up to date
3. **Frontend tests** - Vitest unit tests
4. **Frontend build** - Ensures React app compiles

All checks must pass before merging.

## Architecture Overview

```
┌─────────────────┐     ┌─────────────────┐
│  React Frontend │────▶│   Go Backend    │
│  (Vite + TS)    │     │   (Gin + GORM)  │
└─────────────────┘     └────────┬────────┘
                                 │
                        ┌────────▼────────┐
                        │    SQLite/      │
                        │   PostgreSQL    │
                        └─────────────────┘
```

### Key Technologies

| Layer | Technology |
|-------|------------|
| Frontend | React 19, TypeScript, Vite |
| Backend | Go 1.25, Gin, GORM |
| Database | SQLite (dev), PostgreSQL (prod) |
| Auth | JWT, OIDC/SSO, SCIM 2.0 |
| Docs | Swagger/OpenAPI |
