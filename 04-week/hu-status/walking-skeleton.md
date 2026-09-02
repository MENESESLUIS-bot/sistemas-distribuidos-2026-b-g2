# Walking Skeleton — `library-api`

## Purpose

The walking skeleton proves that one vertical slice can run through the hexagonal layers against a real PostgreSQL database in a container.

## Folder layout

```text
backend/
  cmd/api/main.go                 # composition root
  internal/
    domain/
      catalog/                    # Book aggregate and invariants
    application/
      catalog/                    # use cases and ports
    infrastructure/
      http/                       # handlers and routing
      postgres/                   # adapters and migrations
  migrations/
```

## Endpoints

| Method | Path | Behavior |
|---|---|---|
| GET | `/health` | Returns HTTP 200 and `{"status":"ok"}` when the process is alive. |
| POST | `/api/v1/books` | Persists a book with title, author, ISBN, category, year, and copy count. |
| GET | `/api/v1/books/{id}` | Reads a persisted book from PostgreSQL. |

## Composition root

`cmd/api/main.go` loads environment configuration, creates the PostgreSQL pool, constructs repositories and use cases, registers HTTP handlers, and starts the server. The domain package is constructed without database or HTTP dependencies.

## Verification checklist

- Start PostgreSQL with Docker Compose.
- Apply the migration that creates `books`.
- Start `library-api` with environment variables.
- Call `/health`.
- Create and retrieve a book through the API.
- Restart the API and retrieve the same book to prove persistence.
