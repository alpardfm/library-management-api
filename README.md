# Library Management API

REST API for managing users, books, borrowing, and returns in a small library system.

## What This Project Includes

- JWT authentication with role-based access for `admin`, `librarian`, and `member`
- Book CRUD with search, sorting, and pagination
- Borrow and return flow with transaction boundary in the service layer
- Stock and active-borrow invariants enforced in PostgreSQL
- Unit, integration, and E2E test layers
- GitHub Actions quality gate for lint and unit tests

## Stack

- Go
- Gin
- GORM
- PostgreSQL
- Testify + SQLMock
- golangci-lint

## Quick Start

### 1. Prepare environment

```bash
cp .env.example .env
```

### 2. Start PostgreSQL

```bash
make docker-up
```

### 3. Run the API

```bash
make run
```

API base URL:

```text
http://localhost:8080
```

Health endpoints:

```text
GET /health
GET /ready
```

## Environment Variables

| Variable | Default | Notes |
| --- | --- | --- |
| `APP_PORT` | `8080` | API port |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | PostgreSQL user |
| `DB_PASSWORD` | `password` | PostgreSQL password |
| `DB_NAME` | `library_db` | PostgreSQL database |
| `DB_SSLMODE` | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | `your-super-secret-jwt-key-change-in-production` | JWT signing secret |
| `JWT_EXPIRY` | `24h` | Token expiry |
| `READ_TIMEOUT` | `10s` | HTTP read timeout |
| `WRITE_TIMEOUT` | `10s` | HTTP write timeout |
| `IDLE_TIMEOUT` | `60s` | HTTP idle timeout |
| `MAX_BOOKS_PER_USER` | `5` | Borrow limit per user |
| `BORROW_DAYS` | `14` | Default due date offset |
| `FINE_PER_DAY` | `1000` | Overdue fine per day |

## API Endpoints

### Public

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/api/v1/auth/register` | Register a user |
| `POST` | `/api/v1/auth/login` | Login and get JWT |
| `GET` | `/health` | Liveness check |
| `GET` | `/ready` | Readiness check with DB ping |

### Protected

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/books` | List books |
| `GET` | `/api/v1/books/:id` | Get book detail |
| `POST` | `/api/v1/books` | Create book (`admin`, `librarian`) |
| `PUT` | `/api/v1/books/:id` | Update book (`admin`, `librarian`) |
| `DELETE` | `/api/v1/books/:id` | Delete book (`admin`, `librarian`) |
| `POST` | `/api/v1/borrow` | Borrow a book |
| `POST` | `/api/v1/borrow/return` | Return a book |
| `GET` | `/api/v1/borrow/my-books` | List current user borrows |
| `GET` | `/api/v1/borrow/active` | List active borrows (`admin`, `librarian`) |
| `GET` | `/api/v1/borrow/overdue` | List overdue borrows (`admin`, `librarian`) |

## Response Contract

Success shape:

```json
{
  "success": true,
  "message": "optional message",
  "data": {},
  "meta": {}
}
```

Error shape:

```json
{
  "success": false,
  "message": "book not found",
  "error": {
    "code": "not_found",
    "message": "book not found"
  }
}
```

List query params:

| Param | Description |
| --- | --- |
| `page` | Positive integer, default `1` |
| `limit` | Positive integer, clamped per endpoint |
| `search` | Optional search string |
| `sort` | Optional endpoint-specific sort key |

List meta:

```json
{
  "page": 1,
  "limit": 10,
  "total": 42,
  "total_pages": 5
}
```

## Makefile Commands

| Command | Description |
| --- | --- |
| `make run` | Run API locally |
| `make build` | Build binary |
| `make test` | Run unit and integration tests |
| `make test-unit` | Run unit tests |
| `make test-integration` | Run integration tests |
| `make test-e2e` | Run E2E tests |
| `make lint` | Run golangci-lint |
| `make vet` | Run `go vet` |
| `make quality` | Run lint, vet, and unit tests |
| `make docker-up` | Start PostgreSQL and pgAdmin |
| `make docker-down` | Stop Docker services |

## Test Matrix

| Layer | Command | Purpose |
| --- | --- | --- |
| Unit | `go test ./tests/unit/... -v` | Fast feedback for handlers, services, repositories, middleware |
| Integration | `go test ./tests/integration/... -v` | DB-backed behavior and concurrency checks |
| E2E | `go test ./tests/e2e/... -v` | Happy path against a running API |

### E2E Preconditions

- API server must already be running
- Database must already be available
- Default E2E base URL is `http://localhost:8080`
- Override with `E2E_BASE_URL=http://host:port` if needed

## CI

GitHub Actions workflow:

```text
.github/workflows/quality-gates.yml
```

Required PR jobs:

- `lint`
- `unit-test`

Optional manual job:

- `integration-test`

## ADRs

- [ADR 0001](docs/adr/0001-service-transaction-boundary.md)
- [ADR 0002](docs/adr/0002-database-invariants.md)
- [ADR 0003](docs/adr/0003-standard-response-contract.md)

## Release

- Changelog: [CHANGELOG.md](CHANGELOG.md)
- Release checklist: [RELEASE.md](RELEASE.md)

## Notes

- PostgreSQL-specific constraints and indexes are applied only when the dialector is PostgreSQL.
- `pg_trgm` is enabled gracefully. If extension creation fails, the app continues without trigram indexes.
- Integration concurrency test reference:
  [tests/integration/borrow_concurrency_test.go](https://github.com/alpardfm/library-management-api/blob/master/tests/integration/borrow_concurrency_test.go)
