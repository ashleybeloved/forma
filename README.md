# Forma

[![RU](https://img.shields.io/badge/lang-RU-blue)](./docs/README.ru.md)
[![EN Docs](https://img.shields.io/badge/docs-EN-lightgrey)](./docs/README.en.md)

Forma is a small poll service with a Go backend and a Next.js frontend.

## Stack

- Backend: Go, Gin, SQLite, JWT
- Frontend: Next.js, React, TypeScript, Tailwind CSS

## Features

- user registration and login
- poll creation with multiple questions
- question types:
  - single choice
  - multiple choice
  - text answer
- public poll page by short link
- vote protection by IP and guest token
- auth-only polls
- poll statistics for the author
- vote geography via GeoIP

## Project structure

```text
cmd/server/      backend entry point
internal/        backend code
migrations/      SQL migrations
frontend/        Next.js frontend
data/            SQLite and GeoIP data
```

## Run locally

### Backend

Create `.env` in the project root:

```env
SERVER_PORT=:8080
APP_VERSION=dev
JWT_SECRET_KEY=change_me
DB_PATH=./data/forma.db
GEOIP_DB_PATH=./data/GeoLite2-Country.mmdb
DOMAIN=localhost
HTTPS=false
```

Then run:

```bash
go run ./cmd/server
```

### Frontend

```bash
npm --prefix frontend install
npm --prefix frontend run dev
```

By default:
- backend: `http://localhost:8080`
- frontend: `http://localhost:3000`

## Notes

- backend routes use the `/api` prefix
- frontend talks to backend through the Next.js `/api` proxy
- backend should be started from the repository root
- GeoIP database file is required: `./data/GeoLite2-Country.mmdb`

## Validation

```bash
go test ./...
go build ./...
npm --prefix frontend run build
```

## More docs

- English: `docs/README.en.md`
- Russian: `docs/README.ru.md`
