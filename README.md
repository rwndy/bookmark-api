# Bookmark API

A RESTful API for saving and managing bookmarks/links, built with Go following clean architecture principles.

## Tech Stack

- **Go** — Language
- **Fiber** — HTTP framework
- **GORM** — ORM
- **PostgreSQL** — Database
- **JWT** — Authentication
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

| Method | Endpoint          | Description        | Auth |
|--------|-------------------|--------------------|------|
| POST   | `/auth/register`  | Register new user  | No   |
| POST   | `/auth/login`     | Login, get JWT     | No   |

### Bookmarks

| Method | Endpoint              | Description         | Auth |
|--------|-----------------------|---------------------|------|
| GET    | `/api/bookmarks`      | List my bookmarks   | Yes  |
| POST   | `/api/bookmarks`      | Create bookmark     | Yes  |
| PUT    | `/api/bookmarks/:id`  | Update bookmark     | Yes  |
| DELETE | `/api/bookmarks/:id`  | Delete bookmark     | Yes  |

All protected endpoints require `Authorization: Bearer <token>` header.

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

### Create Bookmark

```bash
curl -X POST http://localhost:3000/api/bookmarks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"url":"https://go.dev","title":"Go Official","tags":"golang,docs"}'
```

### List Bookmarks

```bash
curl http://localhost:3000/api/bookmarks \
  -H "Authorization: Bearer YOUR_TOKEN"
```

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
- API responses follow `{"success": bool, "data": ..., "error": "..."}` format
- Passwords hashed with bcrypt, never stored in plain text
- JWT tokens expire after 24 hours

## License

MIT