# No-as-a-Service

> A lightweight REST API that returns random, creative rejection reasons. Perfect for when you need a polite (or witty) way to say "no".

[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-1.11-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Buy Me A Coffee](https://img.shields.io/badge/Buy_Me_A_Coffee-FFDD00?logo=buy-me-a-coffee&logoColor=black)](https://buymeacoffee.com/martin.willig)


Mainly inspired by https://github.com/hotheadhacker/no-as-a-service

---

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Docker](#docker)
- [Building](#building)
- [Project Structure](#project-structure)
- [License](#license)

---

## Features

- 🎲 Random rejection reasons in English and German
- ⚡ Fast and lightweight
- 🛡️ Built-in rate limiting
- 🔧 Configurable via environment variables
- 🏗️ Cross-platform builds (Linux, macOS, Windows)

---

## Quick Start

### Prerequisites

- Go 1.25 or higher

### Run Development Server

```bash
# Clone the repository
git clone https://github.com/G33kM4sT3r/no-as-a-service.git
cd no-as-a-service

# Install dependencies
make deps

# Start the server
make run
```

The server starts at `http://localhost:8080`

### Test the API

```bash
curl http://localhost:8080/reason
```

---

## Configuration

Create a `.env` file in the project root (or set environment variables):

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DEFAULT_LANGUAGE` | Default language for reasons (`en`, `de`) | `en` |
| `RATE_LIMIT_MAX` | Max requests per window | `120` |
| `RATE_LIMIT_WINDOW_SECONDS` | Rate limit window in seconds | `60` |

**Example `.env`:**

```env
PORT=8080
DEFAULT_LANGUAGE=en
RATE_LIMIT_MAX=120
RATE_LIMIT_WINDOW_SECONDS=60
```

---

## API Reference

### Base URL

```
http://localhost:8080
```

---

### GET `/reason`

Returns a random rejection reason.

#### Query Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `lang` | string | No | Language code: `en` (English) or `de` (German). Defaults to `DEFAULT_LANGUAGE` env variable. |

#### Example Request

```bash
# English (default)
curl http://localhost:8080/reason

# German
curl http://localhost:8080/reason?lang=de
```

#### Success Response (200)

```json
{
  "Payload": {
    "reason": "The vibes said no."
  },
  "Status": {
    "Status": 200,
    "Message": "Reason delivered successfully.",
    "Code": "RE_200_REASON_DELIVERED"
  }
}
```

#### Error Response (400) - Invalid Language

```json
{
  "Payload": null,
  "Status": {
    "Status": 400,
    "Message": "Invalid language. Supported: en, de",
    "Code": "RE_400_INVALID_LANGUAGE"
  }
}
```

---

### GET `/healthcheck`

Health check endpoint for monitoring.

#### Example Request

```bash
curl http://localhost:8080/healthcheck
```

#### Success Response (200)

```json
{
  "Payload": null,
  "Status": {
    "Status": 200,
    "Message": "Service is healthy",
    "Code": "HC_200_OK"
  }
}
```

---

### Rate Limiting

All endpoints are rate-limited by IP address.

- **Default:** 120 requests per 60 seconds
- **Response when exceeded:** `429 Too Many Requests`

```json
{
  "Status": {
    "Status": 429,
    "Message": "Rate limit exceeded. Try again later.",
    "Code": "RL_429_TOO_MANY_REQUESTS"
  }
}
```

---

## Docker

Run the application in a container using Docker or Docker Compose.

### Using Docker Compose (Recommended)

The easiest way to run the application:

```bash
# Build and start the container
docker compose up -d

# View logs
docker compose logs -f

# Stop the container
docker compose down
```

The service will be available at `http://localhost:8080`

#### Environment Variables

You can configure the container using environment variables in your `.env` file or by passing them directly:

```bash
# Using .env file (automatically loaded)
docker compose up -d

# Or override variables directly
PORT=3000 docker compose up -d
```

### Using Docker Directly

#### Build the Image

```bash
docker build -t noaas .
```

#### Run the Container

```bash
docker run -d \
  --name noaas \
  -p 8080:8080 \
  -e PORT=8080 \
  -e DEFAULT_LANGUAGE=en \
  -e RATE_LIMIT_MAX=120 \
  -e RATE_LIMIT_WINDOW_SECONDS=60 \
  noaas
```

### Health Checks

The container includes a health check that monitors the `/healthcheck` endpoint:

```bash
# Check container health status
docker inspect --format='{{.State.Health.Status}}' noaas
```

### Image Details

- **Base Image:** `alpine:3.23` (minimal ~5MB)
- **Multi-stage Build:** Uses `golang:1.25.5-alpine` for building
- **Final Size:** ~15MB (including data files)

---

## Building

### Build for Current Platform

```bash
make build
```

Output: `dist/NoAAS`

### Build for All Platforms

```bash
make build-all
```

Output in `dist/`:
- `NoAAS-linux-amd64`
- `NoAAS-windows-amd64.exe`
- `NoAAS-darwin-arm64` (macOS Apple Silicon)
- `NoAAS-darwin-amd64` (macOS Intel)

### Build for Specific Platform

```bash
make build-linux        # Linux (amd64)
make build-windows      # Windows (amd64)
make build-darwin-arm64 # macOS ARM64
make build-darwin-amd64 # macOS Intel
```

### Other Make Targets

| Command | Description |
|---------|-------------|
| `make run` | Run development server |
| `make test` | Run tests |
| `make deps` | Download dependencies |
| `make lint` | Run linter |
| `make clean` | Remove build artifacts |
| `make help` | Show all available targets |

---

## Project Structure

```
no-as-a-service/
├── main.go                 # Application entry point
├── Makefile                # Build automation
├── Dockerfile              # Container build instructions
├── docker-compose.yml      # Docker Compose configuration
├── .env.sample             # Example environment file
├── data/
│   ├── reasons.en.json     # English rejection reasons
│   └── reasons.de.json     # German rejection reasons
└── internal/
    ├── handler/            # HTTP request handlers
    │   ├── healthcheck.go
    │   └── reason.go
    ├── helper/             # Utility functions
    │   └── env.go
    ├── middleware/         # Gin middleware
    │   └── ratelimit.go
    ├── response/           # Response structures
    │   ├── default.go
    │   └── status.go
    └── router/             # Route definitions
        └── router.go
```

---

## License

MIT License - do whatever, just don't say yes when you should say no.