# Forma

[![RU](https://img.shields.io/badge/lang-RU-blue)](./README.ru.md)
[![Root README](https://img.shields.io/badge/readme-root-black)](../README.md)

Forma is a small poll service with a Go backend and a Next.js frontend.

## Stack

- Backend: Go, Gin, SQLite, JWT
- Frontend: Next.js, React, TypeScript, Tailwind CSS

## Features

- user registration and login
- poll creation with multiple questions
- single choice, multiple choice, and text questions
- public poll page by short link
- vote protection by IP and guest token
- auth-only polls
- statistics and vote geography

## Structure

```text
cmd/server/      backend entry point
internal/        backend code
migrations/      SQL migrations
frontend/        Next.js frontend
data/            SQLite and GeoIP data
```

## Local run

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

Run backend:

```bash
go run ./cmd/server
```

Run frontend:

```bash
npm --prefix frontend install
npm --prefix frontend run dev
```

Default addresses:
- backend: `http://localhost:8080`
- frontend: `http://localhost:3000`

## Notes

- backend routes use `/api`
- frontend uses the Next.js `/api` proxy for backend requests
- run backend from the repository root
- required GeoIP file: `./data/GeoLite2-Country.mmdb`

## Validation

```bash
go test ./...
go build ./...
npm --prefix frontend run build
```
