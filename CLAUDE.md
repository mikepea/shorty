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
- Each folder contains README files explaining the concepts

## CI/CD

- GitHub Actions runs on all PRs and pushes to main
- Tests must pass before merging
- CI runs: unit tests, integration tests, Go builds, and React build
