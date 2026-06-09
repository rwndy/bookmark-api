# Bookmark API

A RESTful API for saving and managing bookmarks/links, built with Go following clean architecture principles.

## Features

- Email + password registration with bcrypt-hashed credentials
- JWT access tokens (short-lived) + opaque refresh tokens (long-lived, SHA-256 hashed at rest)
- Refresh token rotation on every refresh, with revoke-on-logout
- Per-user bookmark CRUD with tag support
- Standardized `{data, message, status}` response envelope
- Database migrations via `golang-migrate`
- Request validation via `go-playground/validator`

## Tech Stack

- **Go** — Language
- **Fiber** — HTTP framework
- **GORM** — ORM
- **PostgreSQL** — Database
- **JWT (HS256)** — Access tokens
- **golang-migrate** — Database migrations
- **go-playground/validator** — Request validation

## Architecture

```
cmd/api/              → Entry point
internal/
├── config/           → Environment configuration
├── domain/           → Structs, interfaces, custom errors
├── handler/          → HTTP handlers (controllers)
│   └── dto/          → Request/response validation
├── service/          → Business logic
├── repository/       → Database queries
└── middleware/       → JWT auth middleware
pkg/response/         → Standardized API responses
migrations/           → SQL migration files
```

Each layer communicates through interfaces. Dependencies flow inward: `handler → service → repository → database`.

## Prerequisites

- Go 1.22+
- PostgreSQL 16+
- [golang-migrate](https://github.com/golang-migrate/migrate) (for database migrations)

## Getting Started

### 1. Clone

```bash
git clone https://github.com/USERNAME/bookmark-api.git
cd bookmark-api
```

### 2. Environment

```bash
cp .env.example .env
```

Edit `.env`:

```env
DB_DSN=postgres://postgres:yourpassword@localhost:5432/bookmark_db?sslmode=disable
JWT_SECRET=your-secret-key
PORT=3000

# Optional — token lifetimes (defaults shown)
ACCESS_TOKEN_TTL_MIN=15
REFRESH_TOKEN_TTL_HOUR=168
```

### 3. Database

```bash
# Create database
psql -U postgres -c "CREATE DATABASE bookmark_db;"

# Run migrations
make migrate-up
```

### 4. Run

```bash
make run
```

Server starts at `http://localhost:3000`.

## API Endpoints

### Auth

| Method | Endpoint          | Description                          | Auth |
|--------|-------------------|--------------------------------------|------|
| POST   | `/auth/register`  | Register new user                    | No   |
| POST   | `/auth/login`     | Login, get access + refresh tokens   | No   |
| POST   | `/auth/refresh`   | Rotate refresh token, get new pair   | No   |
| POST   | `/auth/logout`    | Revoke refresh token                 | No   |

### Bookmarks

| Method | Endpoint              | Description         | Auth |
|--------|-----------------------|---------------------|------|
| GET    | `/api/bookmarks`      | List my bookmarks   | Yes  |
| POST   | `/api/bookmarks`      | Create bookmark     | Yes  |
| PUT    | `/api/bookmarks/:id`  | Update bookmark     | Yes  |
| DELETE | `/api/bookmarks/:id`  | Delete bookmark     | Yes  |

All protected endpoints require `Authorization: Bearer <access_token>` header.

## Authentication Flow

1. **Register** → create an account.
2. **Login** → receive `{ access_token, refresh_token, expires_in }`.
3. Send `Authorization: Bearer <access_token>` with each protected request.
4. When the access token expires, call **`/auth/refresh`** with the refresh token to get a new pair. The previous refresh token is revoked (rotation) — store the new one.
5. **Logout** revokes the refresh token; the corresponding access token remains valid until it expires.

Refresh tokens are random 32-byte hex strings; only their SHA-256 hash is persisted, so a database leak cannot replay them.

## Usage

### Register

```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### Login

```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response:

```json
{
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "a1b2c3...",
    "expires_in": 900
  },
  "message": "login successful",
  "status": 200
}
```

### Refresh

```bash
curl -X POST http://localhost:3000/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### Logout

```bash
curl -X POST http://localhost:3000/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"YOUR_REFRESH_TOKEN"}'
```

### Create Bookmark

```bash
curl -X POST http://localhost:3000/api/bookmarks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{"url":"https://go.dev","title":"Go Official","tags":"golang,docs"}'
```

### List Bookmarks

```bash
curl http://localhost:3000/api/bookmarks \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Response Envelope

All responses share the same shape:

```json
{
  "data": {},
  "message": "human-readable description",
  "status": 200
}
```

`data` is `null` on failure responses.

## Makefile

```bash
make run                             # Start API server
make migrate-up                      # Apply all migrations
make migrate-down                    # Rollback all migrations
make migrate-create name=add_column  # Create new migration
```

## API Spec

Full OpenAPI 3.0 spec available at [`apispec.json`](./apispec.json). Preview it at [editor.swagger.io](https://editor.swagger.io).

## Project Conventions

- `internal/` — Private packages, enforced by Go compiler
- `pkg/` — Public shared utilities
- All errors use `domain.AppError` for consistent HTTP status mapping
- API responses follow the `{data, message, status}` envelope
- Passwords hashed with bcrypt, never stored in plain text
- Access tokens default to 15 min, refresh tokens to 7 days — both configurable via env
- Refresh tokens are opaque, hashed before persistence, and rotated on every refresh
