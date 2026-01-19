# Shorty

A modern, self-hosted URL shortener with team collaboration features, SSO support, and enterprise-ready SCIM provisioning.

## Features

- **URL Shortening** - Create short, memorable links with custom slugs
- **Team Collaboration** - Organize links into groups with role-based access control
- **Tagging System** - Categorize and filter links with tags
- **SSO/OIDC Support** - Integrate with Okta, Azure AD, Keycloak, or any OIDC provider
- **SCIM 2.0 Provisioning** - Automatic user and group sync from your identity provider
- **API Keys** - Programmatic access for automation and integrations
- **Import/Export** - Bulk operations via JSON
- **Admin Dashboard** - User management and system statistics

## Quick Start

### Prerequisites

- Go 1.25+
- Node.js 20+
- SQLite (default) or PostgreSQL

### Running Locally

1. **Clone the repository**
   ```bash
   git clone https://github.com/mikepea/shorty.git
   cd shorty
   ```

2. **Start the backend**
   ```bash
   go run ./cmd/shorty-server
   ```
   The API server starts on `http://localhost:8080`

3. **Start the frontend** (in a separate terminal)
   ```bash
   cd web
   npm install
   npm run dev
   ```
   The frontend starts on `http://localhost:5173`

4. **Login with default admin**
   - Email: `admin@shorty.local`
   - Password: `changeme`

   > **Important:** Change the admin password immediately after first login.

## Configuration

Configure Shorty using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | Database connection string | `shorty.db` (SQLite) |
| `JWT_SECRET` | Secret for JWT tokens | Auto-generated |
| `SHORTY_BASE_URL` | Public URL for the service | `http://localhost:8080` |

### Database Options

**SQLite (default):**
```bash
DATABASE_URL=shorty.db
```

**PostgreSQL:**
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/shorty?sslmode=disable
```

## API Overview

All API endpoints are under `/api`. Authentication is via JWT token in the `Authorization: Bearer <token>` header.

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth/login` | Login |
| `POST` | `/api/auth/register` | Register new user |
| `GET` | `/api/links` | List links |
| `POST` | `/api/links` | Create link |
| `GET` | `/api/groups` | List groups |
| `POST` | `/api/groups` | Create group |
| `GET` | `/api/tags` | List tags |

### SCIM Endpoints

SCIM endpoints are under `/scim/v2` and require a SCIM bearer token.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/scim/v2/Users` | List users |
| `POST` | `/scim/v2/Users` | Create user |
| `GET` | `/scim/v2/Groups` | List groups |
| `POST` | `/scim/v2/Groups` | Create group |
| `GET` | `/scim/v2/ServiceProviderConfig` | SCIM configuration |

## Project Structure

```
shorty/
├── cmd/
│   ├── shorty-cli/        # CLI tool (future)
│   └── shorty-server/     # REST API server
├── pkg/shorty/
│   ├── admin/             # Admin endpoints
│   ├── apikeys/           # API key management
│   ├── auth/              # Authentication
│   ├── groups/            # Group management
│   ├── importexport/      # Bulk operations
│   ├── links/             # Link management
│   ├── models/            # Database models
│   ├── oidc/              # OIDC/SSO support
│   ├── redirect/          # URL redirection
│   ├── scim/              # SCIM 2.0 provisioning
│   └── tags/              # Tag management
├── web/                   # React frontend
│   ├── src/
│   │   ├── api/           # API client
│   │   ├── components/    # UI components
│   │   ├── context/       # React context
│   │   └── pages/         # Page components
│   └── package.json
├── tests/
│   └── integration/       # Integration tests
└── go.mod
```

## Development

### Running Tests

```bash
# Unit tests
go test -v ./pkg/...

# Integration tests
go test -v ./tests/integration/...

# Integration tests with Keycloak (requires Docker)
INTEGRATION_TEST_KEYCLOAK=1 go test -v ./tests/integration/...
```

### Building

```bash
# Build server
go build -o shorty-server ./cmd/shorty-server

# Build frontend
cd web && npm run build
```

## Documentation

- [Deployment Guide](docs/deployment.md) - Production deployment
- [OIDC/SCIM Configuration](docs/oidc-scim.md) - Identity provider setup
- [Admin Guide](docs/admin.md) - Administration tasks
- [Developer Guide](docs/developer.md) - Contributing to Shorty
- [API Reference](docs/api.md) - Full API documentation

## License

MIT

## Contributing

Contributions are welcome! Please read the [Developer Guide](docs/developer.md) for details on the development workflow.
