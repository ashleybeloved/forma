# Forma

![Site](/docs/images/screenshot_1.png)

A pet project poll service built to practice web development with Go and Frontend Layer. It uses Docker Compose to run the backend, frontend, and a reverse proxy together locally.

## Tech Stack

*   **Backend:** Go (1.26+), Gin Web Framework, SQLite, JWT, GeoIP2
*   **Frontend:** Next.js (Standalone mode), React, TypeScript, Tailwind CSS
*   **Reverse Proxy:** Caddy 2 (automated Let's Encrypt HTTPS)
*   **Orchestration:** Docker Compose

## Features

*   **Auth:** User registration and login via JWT.
*   **Dynamic Polls:** Support for multiple question types:
    *   Single choice (Radio buttons)
    *   Multiple choice (Checkboxes)
    *   Text answers (Open-ended)
*  **Short Links:** Public access to polls via unique short URLs.
*  **Anti-Fraud:** Vote protection using client IP tracking and unique guest tokens.
*  **Geo-Analytics:** Voter audience geography tracking via GeoIP database.
*  **Metrics:** Detailed response statistics for poll authors.
*  **Access Control:** Option to restrict polls to authenticated users only.

## Project Structure

```text
├── cmd/server/          # Entry point (main.go)
├── internal/            # Core logic (handlers, models, services)
│   ├── config/          # Application configuration
│   ├── handler/         # HTTP request handlers
│   ├── middleware/      # Middleware (auth, logging, etc.)
│   ├── model/           # Data models
│   ├── pkg/             # Shared packages
│   ├── repository/      # Database access layer
│   └── service/         # Business logic
── migrations/          # SQLite schema migrations
── data/                # Persistent storage (SQLite DB + GeoIP mmdb)
├── frontend/            # Next.js application
│   ├── .dockerignore    
│   └── Dockerfile       # Multi-stage build for standalone output
├── router/              # Route definitions
├── docs/images/         # Documentation assets
├── Caddyfile            # Reverse proxy & SSL config
── docker-compose.yml   # Service orchestration
── Dockerfile           # Backend Docker build file
├── .env.example         # Environment variables template
└── .dockerignore        # Root build exclusions
```

## Local Development

### 1. Prerequisites
Install [Docker Desktop](https://www.docker.com/products/docker-desktop/).

### 2. GeoIP Database
The backend requires the GeoIP database to start (this repository already includes). Place the file here:
`./data/GeoLite2-Country.mmdb`

### 3. Environment Variables
Copy the example environment file and fill in your values:

```bash
cp .env.example .env
```

Then edit `.env` with your configuration.

### 4. Run Services

Build and start containers in the background:

```bash
docker-compose up -d --build
```

## Production

### 1. DNS Configuration

Point your domain to the server's public IP by creating A records:

`your-domain.com` ➔ `SERVER_PUBLIC_IP`

`www.your-domain.com` ➔ `SERVER_PUBLIC_IP`

### 2. Update Caddyfile

Replace localhost with your domain in ./Caddyfile:

```bash
your-domain.com, www.your-domain.com {
    # API traffic -> Go backend
    reverse_proxy /api/* backend:8080

    # Everything else -> Next.js frontend
    reverse_proxy frontend:3000
}
```

### 3. Deploy on your server

`docker-compose up -d --build`

###### Notes: The frontend is completely vibe-coded, but the backend is hand-written. The project may contain mistakes =)
