# Admin Guide

This guide covers administration tasks for Shorty.

## Table of Contents

- [Default Admin Account](#default-admin-account)
- [User Management](#user-management)
- [Group Management](#group-management)
- [SCIM Token Management](#scim-token-management)
- [OIDC Provider Management](#oidc-provider-management)
- [System Statistics](#system-statistics)
- [API Keys](#api-keys)

## Default Admin Account

On first startup, Shorty creates a default admin account:

- **Email**: `admin@shorty.local`
- **Password**: `changeme`

**Important**: Change this password immediately after first login.

### Changing the Admin Password

1. Log in with the default credentials
2. Navigate to Settings
3. Update your password

Or via API:

```bash
# First, login to get a token
TOKEN=$(curl -s -X POST https://your-domain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@shorty.local","password":"changeme"}' | jq -r '.token')

# Then update via the settings endpoint (if available)
# Or create a new admin user and delete the default one
```

## User Management

Admins can view and manage all users in the system.

### Viewing Users

**Via Admin Dashboard**: Navigate to Admin > Users

**Via API**:

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/users
```

Response:
```json
{
  "users": [
    {
      "id": 1,
      "email": "admin@shorty.local",
      "name": "Admin",
      "system_role": "admin",
      "active": true,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 1
}
```

### User Roles

Shorty has two system-level roles:

| Role | Description |
|------|-------------|
| `user` | Standard user, can create links and groups |
| `admin` | Full system access, can manage users and settings |

### Deactivating Users

Deactivated users cannot log in but their data is preserved.

```bash
curl -X PATCH \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/admin/users/123 \
  -d '{"active": false}'
```

### Promoting a User to Admin

```bash
curl -X PATCH \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/admin/users/123 \
  -d '{"system_role": "admin"}'
```

## Group Management

Groups organize links and control access. Each group has its own members with roles.

### Group Roles

| Role | Permissions |
|------|-------------|
| `viewer` | View links in the group |
| `editor` | View, create, edit, delete links |
| `admin` | All editor permissions + manage members |

### Viewing All Groups

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/groups
```

### Group Statistics

View link counts and member counts for all groups:

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/stats
```

## SCIM Token Management

SCIM tokens authenticate identity providers for user/group provisioning.

### Creating a SCIM Token

**Via Admin Dashboard**: Admin > SCIM Tokens > Create Token

**Via API**:

```bash
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/admin/scim-tokens \
  -d '{"description": "Okta Production"}'
```

Response:
```json
{
  "id": 1,
  "token": "abc123def456...",
  "token_prefix": "abc123de",
  "description": "Okta Production",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Important**: The full token is only shown once. Store it securely.

### Listing SCIM Tokens

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/scim-tokens
```

Response shows token prefixes (not full tokens) and last used timestamps:

```json
{
  "tokens": [
    {
      "id": 1,
      "token_prefix": "abc123de",
      "description": "Okta Production",
      "last_used_at": "2024-01-15T12:00:00Z",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

### Revoking a SCIM Token

```bash
curl -X DELETE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/scim-tokens/1
```

## OIDC Provider Management

Manage SSO identity providers.

### Adding an OIDC Provider

```bash
curl -X POST \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/admin/oidc/providers \
  -d '{
    "name": "Okta",
    "slug": "okta",
    "issuer": "https://your-org.okta.com",
    "client_id": "your-client-id",
    "client_secret": "your-client-secret",
    "scopes": "openid profile email",
    "enabled": true,
    "auto_provision": true
  }'
```

### Listing OIDC Providers

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/oidc/providers
```

### Disabling an OIDC Provider

```bash
curl -X PATCH \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/admin/oidc/providers/1 \
  -d '{"enabled": false}'
```

### Deleting an OIDC Provider

```bash
curl -X DELETE \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/oidc/providers/1
```

## System Statistics

View overall system statistics.

### Getting System Stats

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  https://your-domain.com/api/admin/stats
```

Response:
```json
{
  "users": {
    "total": 150,
    "active": 142,
    "admins": 3
  },
  "groups": {
    "total": 45
  },
  "links": {
    "total": 1234,
    "public": 89
  }
}
```

## API Keys

Users can create API keys for programmatic access. As an admin, you can view API key usage across the system.

### How API Keys Work

- API keys provide the same access as the user who created them
- Keys can be used instead of JWT tokens for API calls
- Keys don't expire but can be revoked

### User API Key Management

Users manage their own keys at Settings > API Keys or via:

```bash
# List user's API keys
curl -H "Authorization: Bearer $USER_TOKEN" \
  https://your-domain.com/api/api-keys

# Create a new API key
curl -X POST \
  -H "Authorization: Bearer $USER_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-domain.com/api/api-keys \
  -d '{"name": "CI/CD Integration"}'

# Delete an API key
curl -X DELETE \
  -H "Authorization: Bearer $USER_TOKEN" \
  https://your-domain.com/api/api-keys/1
```

### Using API Keys

API keys are used in the Authorization header, same as JWT tokens:

```bash
curl -H "Authorization: Bearer $API_KEY" \
  https://your-domain.com/api/links
```

## Monitoring

### Health Check

Shorty exposes a health endpoint for monitoring:

```bash
curl https://your-domain.com/health
# {"status":"ok"}

curl https://your-domain.com/api/health
# {"status":"ok","service":"shorty"}
```

### Recommended Monitoring

- Monitor the `/health` endpoint for uptime
- Track response times for the redirect endpoint (`/{slug}`)
- Monitor database connection pool usage
- Set up alerts for error rates in logs

## Backup and Recovery

### Database Backup

For PostgreSQL:

```bash
# Backup
pg_dump -U shorty shorty > backup-$(date +%Y%m%d).sql

# Restore
psql -U shorty shorty < backup-20240115.sql
```

### What to Back Up

1. **Database**: Contains all users, groups, links, and configuration
2. **Environment variables**: JWT_SECRET is critical (losing it invalidates all tokens)

### Disaster Recovery

1. Restore the database from backup
2. Ensure environment variables are set (especially JWT_SECRET)
3. Start the application
4. Verify health endpoint responds
5. Test login functionality
