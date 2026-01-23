# Claude Code Instructions

## Git Workflow

- Always use feature branches for changes (e.g., `feature/description`)
- Create separate commits for each logical step
- Push branches and create PRs - never commit directly to main
- Branch protection is enabled: PRs require passing CI checks before merge
- Use `gh pr create` and `gh pr merge` for PR operations

## Commit Message Format

Use this format for all commits:
```
Short description of change

Optional longer description if needed.

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
```

Use heredoc syntax to preserve formatting:
```bash
git commit -m "$(cat <<'EOF'
Commit message here

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude Opus 4.5 <noreply@anthropic.com>
EOF
)"
```

## Project Structure

This is a Go project with library/CLI separation and a React frontend:

```
├── cmd/
│   ├── shorty-cli/     # CLI tool
│   └── shorty-server/  # REST API server
├── pkg/shorty/          # Library package (reusable)
├── web/                # React frontend (Vite + React 18)
│   ├── src/
│   │   ├── components/ # UI components
│   │   ├── hooks/      # Custom React hooks
│   │   └── api/        # API client layer
│   └── package.json
├── go.mod              # Go 1.25
└── .github/workflows/  # CI configuration
```

## Go

- Use Go 1.25
- Library code goes in `pkg/` for reuse
- CLI tools and servers go in `cmd/`
- Run tests with `go test ./... -v`
- Build CLI: `go build ./cmd/hn-scraper`
- Build server: `go build ./cmd/hn-server`

## REST API


## React Frontend

- Built with Vite and React 18
- Run dev server: `cd web && npm run dev`
- Build for production: `cd web && npm run build`
- Vite proxies API requests to Go server in development
- Run tests: `cd web && npm test`

### React Beginner Context

**The user is learning React.** When writing or modifying React code:

1. **Add explanatory comments** that explain React concepts, not just what the code does:
   - Explain hooks (useState, useEffect, useContext) when you use them
   - Explain JSX patterns (conditional rendering, mapping, fragments)
   - Explain TypeScript patterns (interfaces, generics, optional chaining)

2. **Structure components clearly** with sections:
   ```tsx
   // ============================================================================
   // State Management
   // ============================================================================

   // ============================================================================
   // Effects
   // ============================================================================

   // ============================================================================
   // Event Handlers
   // ============================================================================

   // ============================================================================
   // Render
   // ============================================================================
   ```

3. **Key patterns to explain when encountered:**
   - Controlled inputs (value + onChange)
   - Conditional rendering (`&&` and ternary)
   - List rendering with `.map()` and keys
   - Async operations in handlers (try/catch/finally)
   - Custom hooks and context
   - Event handling (preventDefault, e.target.value)

4. **Keep comments up to date** when modifying commented code

## CI/CD

- GitHub Actions runs on all PRs and pushes to main
- Tests must pass before merging
- CI runs: unit tests, integration tests, Go builds, and React build

## Documentation

Developer documentation lives in `docs/` and should be kept up to date:

| File | Content |
|------|---------|
| `docs/developer.md` | Index page linking to all dev guides |
| `docs/setup.md` | Development environment setup |
| `docs/backend.md` | Go backend development |
| `docs/frontend.md` | React frontend development |
| `docs/diagnosis.md` | Debugging and troubleshooting |

**When to update docs:**
- Adding new dependencies or tools → update `setup.md`
- Adding new API endpoints or packages → update `backend.md`
- Adding new frontend features or test patterns → update `frontend.md`
- Adding new debugging techniques or common issues → update `diagnosis.md`

Include doc updates in the same PR as the feature when possible.
