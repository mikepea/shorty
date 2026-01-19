# OIDC and SCIM Configuration Guide

This guide covers setting up Single Sign-On (SSO) via OIDC and automated user provisioning via SCIM 2.0.

## Table of Contents

- [Overview](#overview)
- [OIDC Configuration](#oidc-configuration)
  - [Generic OIDC Setup](#generic-oidc-setup)
  - [Okta](#okta)
  - [Azure AD](#azure-ad)
  - [Keycloak](#keycloak)
- [SCIM Configuration](#scim-configuration)
  - [Creating a SCIM Token](#creating-a-scim-token)
  - [Okta SCIM Setup](#okta-scim-setup)
  - [Azure AD SCIM Setup](#azure-ad-scim-setup)
- [Troubleshooting](#troubleshooting)

## Overview

Shorty supports:

- **OIDC (OpenID Connect)**: Allows users to authenticate using their identity provider (IdP) credentials
- **SCIM 2.0**: Enables automatic provisioning and deprovisioning of users and groups

Both features can be used independently or together for a complete identity management solution.

## OIDC Configuration

### Prerequisites

- Admin access to Shorty
- Admin access to your identity provider
- `SHORTY_BASE_URL` environment variable set correctly

### Generic OIDC Setup

1. **In your Identity Provider**, create a new OIDC/OAuth2 application:
   - Application type: Web application
   - Redirect URI: `https://your-shorty-domain.com/api/oidc/callback`
   - Scopes: `openid`, `profile`, `email`

2. **Note the following values** from your IdP:
   - Client ID
   - Client Secret
   - Issuer URL (e.g., `https://your-idp.com/realms/your-realm`)

3. **In Shorty**, navigate to Admin > OIDC Providers and add a new provider:

   ```json
   {
     "name": "My Identity Provider",
     "slug": "my-idp",
     "issuer": "https://your-idp.com/realms/your-realm",
     "client_id": "your-client-id",
     "client_secret": "your-client-secret",
     "scopes": "openid profile email",
     "enabled": true,
     "auto_provision": true
   }
   ```

   Or via API:

   ```bash
   curl -X POST https://your-shorty-domain.com/api/admin/oidc/providers \
     -H "Authorization: Bearer $ADMIN_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "My Identity Provider",
       "slug": "my-idp",
       "issuer": "https://your-idp.com/realms/your-realm",
       "client_id": "your-client-id",
       "client_secret": "your-client-secret",
       "scopes": "openid profile email",
       "enabled": true,
       "auto_provision": true
     }'
   ```

### Okta

1. **In Okta Admin Console**:
   - Go to Applications > Create App Integration
   - Select "OIDC - OpenID Connect" and "Web Application"
   - Set Sign-in redirect URI: `https://your-shorty-domain.com/api/oidc/callback`
   - Set Sign-out redirect URI: `https://your-shorty-domain.com`
   - Under Assignments, assign users/groups who should have access

2. **Copy the credentials**:
   - Client ID (from General tab)
   - Client Secret (from General tab)
   - Issuer URL: `https://your-okta-domain.okta.com`

3. **Add to Shorty**:

   ```json
   {
     "name": "Okta",
     "slug": "okta",
     "issuer": "https://your-okta-domain.okta.com",
     "client_id": "0oaxxxxxxxxxxxxxx",
     "client_secret": "your-client-secret",
     "scopes": "openid profile email",
     "enabled": true,
     "auto_provision": true
   }
   ```

### Azure AD

1. **In Azure Portal**:
   - Go to Azure Active Directory > App registrations > New registration
   - Name: "Shorty"
   - Redirect URI: Web, `https://your-shorty-domain.com/api/oidc/callback`
   - Click Register

2. **Configure the app**:
   - Go to Certificates & secrets > New client secret
   - Copy the secret value immediately (it won't be shown again)
   - Go to API permissions > Add permission > Microsoft Graph > Delegated > `openid`, `profile`, `email`

3. **Copy the credentials**:
   - Application (client) ID (from Overview)
   - Client Secret (from Certificates & secrets)
   - Issuer URL: `https://login.microsoftonline.com/{tenant-id}/v2.0`

4. **Add to Shorty**:

   ```json
   {
     "name": "Azure AD",
     "slug": "azure",
     "issuer": "https://login.microsoftonline.com/your-tenant-id/v2.0",
     "client_id": "your-application-id",
     "client_secret": "your-client-secret",
     "scopes": "openid profile email",
     "enabled": true,
     "auto_provision": true
   }
   ```

### Keycloak

1. **In Keycloak Admin Console**:
   - Select your realm (or create one)
   - Go to Clients > Create client
   - Client ID: `shorty`
   - Client Protocol: `openid-connect`
   - Root URL: `https://your-shorty-domain.com`

2. **Configure the client**:
   - Access Type: `confidential`
   - Valid Redirect URIs: `https://your-shorty-domain.com/api/oidc/callback`
   - Go to Credentials tab and copy the Secret

3. **Add to Shorty**:

   ```json
   {
     "name": "Keycloak",
     "slug": "keycloak",
     "issuer": "https://your-keycloak-domain.com/realms/your-realm",
     "client_id": "shorty",
     "client_secret": "your-client-secret",
     "scopes": "openid profile email",
     "enabled": true,
     "auto_provision": true
   }
   ```

## SCIM Configuration

SCIM (System for Cross-domain Identity Management) enables automatic user and group provisioning from your identity provider.

### SCIM Endpoints

Shorty implements SCIM 2.0 with the following endpoints:

| Endpoint | Description |
|----------|-------------|
| `GET /scim/v2/ServiceProviderConfig` | SCIM configuration |
| `GET /scim/v2/ResourceTypes` | Available resource types |
| `GET /scim/v2/Schemas` | SCIM schemas |
| `GET /scim/v2/Users` | List users |
| `POST /scim/v2/Users` | Create user |
| `GET /scim/v2/Users/{id}` | Get user |
| `PUT /scim/v2/Users/{id}` | Replace user |
| `PATCH /scim/v2/Users/{id}` | Update user |
| `DELETE /scim/v2/Users/{id}` | Delete user |
| `GET /scim/v2/Groups` | List groups |
| `POST /scim/v2/Groups` | Create group |
| `GET /scim/v2/Groups/{id}` | Get group |
| `PATCH /scim/v2/Groups/{id}` | Update group (add/remove members) |
| `DELETE /scim/v2/Groups/{id}` | Delete group |

### Creating a SCIM Token

1. **Via Admin UI**: Navigate to Admin > SCIM Tokens > Create Token

2. **Via API**:

   ```bash
   curl -X POST https://your-shorty-domain.com/api/admin/scim-tokens \
     -H "Authorization: Bearer $ADMIN_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"description": "Okta SCIM Integration"}'
   ```

   Response:
   ```json
   {
     "token": "abc123...",
     "token_prefix": "abc123",
     "description": "Okta SCIM Integration"
   }
   ```

   **Important**: Save the token immediately. It cannot be retrieved later.

### Okta SCIM Setup

1. **In Okta Admin Console**:
   - Go to your Shorty app > Provisioning > Configure API Integration
   - Check "Enable API Integration"
   - SCIM connector base URL: `https://your-shorty-domain.com/scim/v2`
   - Unique identifier field: `userName`
   - Authentication Mode: `HTTP Header`
   - Authorization: `Bearer your-scim-token`
   - Click "Test API Credentials"

2. **Enable provisioning features**:
   - Go to Provisioning > To App
   - Enable: Create Users, Update User Attributes, Deactivate Users
   - Optionally enable: Sync Password (if using password sync)

3. **Configure attribute mappings**:
   - Map Okta attributes to SCIM attributes:
     - `userName` → `email`
     - `givenName` → `name.givenName`
     - `familyName` → `name.familyName`
     - `email` → `emails[primary eq true].value`

4. **Push groups** (optional):
   - Go to Push Groups
   - Add groups to push to Shorty

### Azure AD SCIM Setup

1. **In Azure Portal**:
   - Go to your Shorty Enterprise Application > Provisioning
   - Set Provisioning Mode to "Automatic"
   - Under Admin Credentials:
     - Tenant URL: `https://your-shorty-domain.com/scim/v2`
     - Secret Token: Your SCIM token
   - Click "Test Connection"

2. **Configure mappings**:
   - Go to Mappings > Provision Azure Active Directory Users
   - Ensure these mappings exist:
     - `userPrincipalName` → `userName`
     - `mail` → `emails[type eq "work"].value`
     - `givenName` → `name.givenName`
     - `surname` → `name.familyName`
     - `displayName` → `displayName`

3. **Start provisioning**:
   - Set Provisioning Status to "On"
   - Click Save

## Troubleshooting

### OIDC Issues

**"Invalid redirect URI" error**
- Ensure `SHORTY_BASE_URL` matches your public URL exactly
- Check that the redirect URI in your IdP matches: `{SHORTY_BASE_URL}/api/oidc/callback`

**"Invalid issuer" error**
- Verify the issuer URL is correct (check for trailing slashes)
- For Azure AD, ensure you're using the v2.0 endpoint

**User not provisioned after SSO login**
- Check that `auto_provision` is enabled for the OIDC provider
- Verify the IdP is returning email in the claims

### SCIM Issues

**401 Unauthorized**
- Verify the SCIM token is correct
- Ensure the Authorization header format is: `Bearer {token}`

**User not created**
- Check that `userName` is a valid email address
- Verify required fields are being sent (userName, name)

**Group membership not syncing**
- Ensure the user exists before adding to a group
- Check that the user ID in the membership request is correct

### Testing SCIM Manually

```bash
# Test authentication
curl -H "Authorization: Bearer $SCIM_TOKEN" \
  https://your-shorty-domain.com/scim/v2/ServiceProviderConfig

# List users
curl -H "Authorization: Bearer $SCIM_TOKEN" \
  https://your-shorty-domain.com/scim/v2/Users

# Create a user
curl -X POST \
  -H "Authorization: Bearer $SCIM_TOKEN" \
  -H "Content-Type: application/json" \
  https://your-shorty-domain.com/scim/v2/Users \
  -d '{
    "schemas": ["urn:ietf:params:scim:schemas:core:2.0:User"],
    "userName": "test@example.com",
    "name": {
      "givenName": "Test",
      "familyName": "User"
    },
    "emails": [{"value": "test@example.com", "primary": true}],
    "active": true
  }'
```

### Logs

Enable debug logging for more details:

```bash
# Check server logs for OIDC/SCIM errors
journalctl -u shorty -f | grep -E "(oidc|scim)"
```
