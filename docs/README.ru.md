# Forma
![Site](docs/images/screenshot_1.png)

[![EN](https://img.shields.io/badge/lang-EN-black)](./README.en.md)
[![Root README](https://img.shields.io/badge/readme-root-black)](../README.md)

Forma — небольшой сервис для создания опросов с backend на Go и frontend на Next.js.

## Стек

- Backend: Go, Gin, SQLite, JWT
- Frontend: Next.js, React, TypeScript, Tailwind CSS

## Возможности

- регистрация и вход
- создание опросов из нескольких вопросов
- типы вопросов: один вариант, несколько вариантов, текстовый ответ
- публичная страница опроса по короткой ссылке
- защита голосования по IP и guest token
- опросы только для авторизованных
- статистика и география голосов

## Структура

```text
cmd/server/      точка входа backend
internal/        backend код
migrations/      SQL-миграции
frontend/        Next.js frontend
data/            SQLite и GeoIP данные
```

## Локальный запуск

Создайте `.env` в корне проекта:

```env
SERVER_PORT=:8080
APP_VERSION=dev
JWT_SECRET_KEY=change_me
DB_PATH=./data/forma.db
GEOIP_DB_PATH=./data/GeoLite2-Country.mmdb
DOMAIN=localhost
HTTPS=false
```

Запуск backend:

```bash
go run ./cmd/server
```

Запуск frontend:

```bash
npm --prefix frontend install
npm --prefix frontend run dev
```

Адреса по умолчанию:
- backend: `http://localhost:8080`
- frontend: `http://localhost:3000`

## Важно

- backend маршруты используют префикс `/api`
- frontend ходит в backend через Next.js proxy `/api`
- backend нужно запускать из корня репозитория
- нужен GeoIP файл: `./data/GeoLite2-Country.mmdb`

## Проверка

```bash
go test ./...
go build ./...
npm --prefix frontend run build
```
