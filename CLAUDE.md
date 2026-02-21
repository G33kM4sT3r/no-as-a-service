# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

No-as-a-Service (NoAAS) — a Go REST API that returns random rejection reasons with multi-language support and IP-based rate limiting. Built with the Gin web framework.

## Commands

```bash
make run          # Start development server
make build        # Build binary → dist/NoAAS
make test         # Run tests (go test -v ./...)
make lint         # Run golangci-lint
make deps         # Download and tidy dependencies
make build-all    # Cross-compile for Linux, Windows, macOS (ARM64 + Intel)
make clean        # Remove build artifacts
```

Docker:
```bash
docker compose up -d    # Build and start
docker compose down     # Stop
```

## Architecture

**Entry point:** `main.go` — loads `.env`, configures Gin engine, adds rate limiting middleware, sets up routes, starts server.

**Request flow:** Gin Engine → Rate Limit Middleware → Router → Handler → JSON Response

### Key packages under `internal/`:

- **router/** — Route definitions: `GET /reason`, `GET /healthcheck`, plus 404 handler
- **handler/** — Endpoint logic. `reason.go` is the core: dynamically discovers languages from `data/reasons.*.json` files, caches loaded reasons in memory with `sync.RWMutex`
- **middleware/** — IP-based rate limiter using `sync.Mutex` with per-IP request counting and time-window expiration
- **response/** — Standardized response structs (`DefaultResponse` wrapping `StatusResponse`)
- **helper/** — `GetEnv`/`GetEnvInt` for environment variable reading with fallbacks

### Data

Rejection reasons live in `data/reasons.{lang}.json` as simple JSON string arrays. Adding a new language only requires adding a new file following that naming pattern — the handler discovers languages at runtime by scanning the `data/` directory.

### Response code convention

Status codes follow `{PREFIX}_{HTTP_STATUS}_{DESCRIPTION}` format:
- `RE_` = reason endpoint
- `HC_` = healthcheck
- `RL_` = rate limiter
- `RT_` = router

## Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `PORT` | `8080` | Server port |
| `DEFAULT_LANGUAGE` | `en` | Fallback language for `/reason` |
| `RATE_LIMIT_MAX` | `120` | Max requests per window per IP |
| `RATE_LIMIT_WINDOW_SECONDS` | `60` | Rate limit window in seconds |
