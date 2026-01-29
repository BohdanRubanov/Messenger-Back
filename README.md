# Messenger-Back — Go REST API (Products + Users/Auth)

A lightweight REST API written in Go with a clean separation of **handlers → services → repositories** and PostgreSQL access via **pgxpool**.

## Features

### Products
- Create product
- Get all products
- Get product by ID
- Update product (partial update via `COALESCE`)
- Delete product

### Users / Auth
- Registration with **Argon2id** password hashing + **pepper** (stored in env)
- Authorization (email + password verification)
- Get all users (without password)
- Get user by ID (without password)
- Update user (partial update via `COALESCE`)
- Delete user

> Note: Authorization currently returns `UserWithoutPassword`. Token-based auth (JWT, refresh tokens, etc.) can be added on top.

## Tech stack

- **Go** (net/http)
- **PostgreSQL**
- **pgx v5 / pgxpool**
- **Argon2id** (`golang.org/x/crypto/argon2`)
- **Docker Compose** (local DB)
- **godotenv** (local env loading)

## Project structure

```text
lesson-proj/
├── cmd/api/                  # App entrypoint + HTTP wiring
│   ├── main.go               # Bootstraps DB, services, handlers, routes
│   ├── middlewares.go        # Logging + CORS middleware
│   └── utils.go              # Routing helpers (method handler, id parsing)
├── internal/
│   ├── database/             # Repositories (SQL/pgxpool access)
│   │   ├── database.go       # pgxpool Connect()
│   │   ├── products.go       # ProductRepository
│   │   └── users.go          # UserRepository
│   ├── handlers/             # HTTP handlers (JSON decode/encode)
│   │   ├── auth.go           # UserHandler (register/auth + CRUD)
│   │   ├── product.go        # ProductHandler (CRUD)
│   │   └── utils.go          # respondWithJSON/respondWithError helpers
│   ├── models/               # Request/response models
│   │   ├── product.go
│   │   └── user.go
│   └── services/             # Business logic (validation, hashing)
│       ├── auth/
│       │   ├── auth.go
│       │   └── utils/
│       │       ├── config.go       # Argon2 params constants
│       │       ├── password.go     # HashPassword/VerifyPassword
│       │       └── validation.go   # User input validation
│       └── products/
│           ├── products.go
│           └── utils/
│               └── validation.go   # Product input validation
├── sql/
│   └── init.sql              # Initial schema (executed on first DB init)
├── docker-compose.yml
├── .env
├── .env.example
├── go.mod
└── go.sum
```

## Requirements

- Go installed (project uses Go in `go.mod`)
- Docker + Docker Compose
- PostgreSQL client optional (psql / VS Code extension)

## Environment variables

Create a `.env` file (or copy `.env.example`) and set:

```env
DB_NAME=lesson
DB_USER=lesson_user
DB_PASSWORD=your_password
DATABASE_URL=postgres://lesson_user:your_password@localhost:5432/lesson?sslmode=disable
SERVER_PORT=8080

# Secret pepper (do NOT commit real value)
PASSWORD_PEPPER=change_me_to_a_long_random_secret
```

### Notes on PASSWORD_PEPPER
- Pepper is an **application secret**, not stored in the DB.
- Keep it in env / secret manager.
- If the pepper leaks, attacker can verify guesses faster; treat it like a key.

## Running locally

### 1) Start PostgreSQL with Docker Compose

```bash
docker compose up -d
```

Verify:
```bash
docker ps
```

### 2) Run the API

From the repo root:

```bash
go run ./cmd/api
```

The server listens on:
```text
http://localhost:${SERVER_PORT}
```

## Database schema init (`sql/init.sql`)

The mounted `sql/init.sql` runs **only on first DB creation** (when the volume is empty).  
If you want it to run again, reset the volume:

```bash
docker compose down
docker volume rm lessonproj_postgres_data   # volume name may differ in your setup
docker compose up -d
```

## API endpoints

### Products

- `GET /products` — list products
- `POST /products/create` — create product
- `GET /products/{id}` — get product by ID
- `PUT /products/{id}` — update product (partial)
- `DELETE /products/{id}` — delete product

#### Create product example

```bash
curl -X POST http://localhost:8080/products/create \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Phone",
    "description": "New model",
    "price": 999
  }'
```

### Users / Auth

- `GET /users` — list users (without password)
- `POST /users/create` — register user
- `POST /users/auth` — authorize user (email + password)
- `GET /users/{id}` — get user by ID (without password)
- `PUT /users/{id}` — update user (partial)
- `DELETE /users/{id}` — delete user

#### Registration example

```bash
curl -X POST http://localhost:8080/users/create \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test",
    "password": "strong_password"
  }'
```

#### Authorization example

```bash
curl -X POST http://localhost:8080/users/auth \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "strong_password"
  }'
```

## Password hashing details

- Algorithm: **Argon2id**
- Salt: random per password, stored inside the hash string (Base64)
- Pepper: `PASSWORD_PEPPER` from env
- Stored format (example):

```text
argon2id$v=19$t=2$m=65536$p=2$<salt>$<hash>
```

Storing parameters + salt with the hash is standard practice; it allows verification even if you change defaults later.

## Production notes

- Replace `Access-Control-Allow-Origin: *` with your real frontend origin(s).
- Add rate limiting for registration/auth endpoints.
- Add request body size limits.
- Add structured logging + request IDs.
- Add JWT + refresh tokens (or session store) if you need stateless auth for clients.
- Run migrations via a migration tool (see below).

## Migrations

Right now the project uses `sql/init.sql` for first-time DB init. For production, use a migration tool such as:
- `golang-migrate/migrate`
- `pressly/goose`

A typical approach:
- Create `migrations/` with versioned `*.up.sql` / `*.down.sql`
- Run migrations on deploy / startup (carefully), or in CI/CD.

## License

Private / learning project (add license if needed).
