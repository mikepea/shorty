# Debugging & Diagnostics

This guide covers troubleshooting and debugging techniques for Shorty.

## Server Logs

The server logs to stdout. Useful log messages include:
- Database connection status
- Migration results
- Authentication errors
- Request/response details (in debug mode)

### Log Levels

By default, the server logs at INFO level. Key log messages:

```
INFO: Server starting on :8080
INFO: Database connected
INFO: Running migrations...
INFO: OIDC provider configured: google
```

### Request Logging

Gin's default logger shows all HTTP requests:

```
[GIN] 2024/01/20 - 10:30:45 | 200 |     1.234ms |    127.0.0.1 | GET     "/api/links"
[GIN] 2024/01/20 - 10:30:46 | 401 |     0.456ms |    127.0.0.1 | POST    "/api/auth/login"
```

## Database Debugging

### Enable SQL Query Logging

To see all SQL queries, modify the database connection in `pkg/shorty/database/database.go`:

```go
import "gorm.io/gorm/logger"

db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})
```

This outputs all queries:

```
[INFO] SELECT * FROM users WHERE id = 1
[INFO] INSERT INTO links (slug, url, ...) VALUES (?, ?, ...)
```

### Inspecting SQLite Database

```bash
# Open the database
sqlite3 shorty.db

# List tables
.tables

# Show schema
.schema users
.schema links

# Query data
SELECT * FROM users;
SELECT slug, url, click_count FROM links LIMIT 10;

# Exit
.quit
```

### Common Database Issues

#### "database is locked"

SQLite only allows one writer at a time. This can happen with concurrent requests:

1. Check for long-running transactions
2. Consider using PostgreSQL for production
3. Ensure connections are properly closed

#### Migration Failures

If migrations fail on startup:

1. Check the error message for which migration failed
2. Back up the database
3. Manually apply or skip the problematic migration
4. Restart the server

## API Debugging

### Using curl

```bash
# Login and get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@shorty.local","password":"changeme"}' \
  | jq -r '.token')

# Use token for authenticated requests
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/links

# Create a link
curl -X POST http://localhost:8080/api/groups/1/links \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com","title":"Example"}'
```

### Swagger UI

Interactive API documentation is available at:

```
http://localhost:8080/swagger/index.html
```

### Common API Errors

| Status | Error | Cause |
|--------|-------|-------|
| 401 | "Authentication required" | Missing or invalid JWT token |
| 401 | "Invalid email or password" | Wrong login credentials |
| 403 | "Admin access required" | User lacks admin role |
| 404 | "Link not found" | Invalid slug or no access |
| 409 | "Email already registered" | Duplicate registration |
| 409 | "Slug already exists" | Duplicate link slug |

## Frontend Debugging

### React Developer Tools

Install the [React Developer Tools](https://react.dev/learn/react-developer-tools) browser extension to:

- Inspect component hierarchy
- View component props and state
- Profile render performance

### Network Tab

Use the browser's Network tab to debug API calls:

1. Open DevTools (F12)
2. Go to Network tab
3. Filter by "Fetch/XHR"
4. Inspect request/response details

### Common Frontend Issues

#### "Failed to fetch" or CORS errors

The Vite dev server proxies `/api` requests to the backend. Ensure:

1. Backend is running on `http://localhost:8080`
2. Vite config has correct proxy settings
3. You're accessing the frontend at `http://localhost:3000`

#### Authentication not persisting

The JWT token is stored in `localStorage`. Check:

1. Token is being saved after login
2. Token isn't expired
3. No JavaScript errors preventing storage

```javascript
// In browser console
localStorage.getItem('token')
```

#### State not updating

React state issues can be debugged by:

1. Adding `console.log` in useEffect hooks
2. Using React DevTools to inspect state
3. Checking for missing dependency arrays

## OIDC/SSO Debugging

### Test OIDC Flow

1. Check provider configuration:
   ```bash
   curl http://localhost:8080/api/oidc/providers
   ```

2. Verify the callback URL is correct in your IdP settings:
   ```
   http://localhost:8080/api/oidc/callback/{provider-slug}
   ```

3. Check server logs for OIDC errors during authentication

### Common OIDC Issues

| Issue | Cause | Solution |
|-------|-------|----------|
| "Provider not found" | Invalid slug in URL | Check provider slug in database |
| "Invalid state" | State mismatch | Clear cookies, try again |
| "Token exchange failed" | Wrong client secret | Verify client_secret in provider config |
| "User not provisioned" | auto_provision disabled | Enable auto_provision or pre-create user |

## Performance Issues

### Slow API Responses

1. Enable SQL logging to identify slow queries
2. Check for N+1 query problems (use preloading)
3. Add database indexes for frequently queried columns

### Memory Usage

Monitor Go memory usage:

```go
import "runtime"

var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Alloc = %v MiB", m.Alloc / 1024 / 1024)
```

### Frontend Performance

1. Use React DevTools Profiler
2. Check bundle size with `npm run build`
3. Lazy load routes for large applications

## Getting Help

If you're stuck:

1. Check existing [GitHub Issues](https://github.com/mikepea/shorty/issues)
2. Search error messages in the codebase
3. Review test files for expected behavior
4. Open a new issue with:
   - Steps to reproduce
   - Expected vs actual behavior
   - Relevant logs/screenshots
