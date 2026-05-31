# GoBanking v2 (Phase 1)

## What’s included
- Clean Architecture project skeleton
- Config via environment variables
- Structured logging setup (zap)
- PostgreSQL connection pool (pgx)
- Health endpoint: `GET /health`
- Docker + docker-compose (API + Postgres)
- Migration folder + example migration
- Graceful shutdown (SIGINT/SIGTERM)

## Run locally (Docker)
1. Build and start:
   - `docker compose up --build`
2. Health check:
   - `curl http://localhost:8080/health`

## Migrations
This project expects the `migrate` CLI from https://github.com/golang-migrate/migrate.

- Apply:
  - `sh ./scripts/migrate-up.sh`
- Rollback:
  - `sh ./scripts/migrate-down.sh`

> Note: Phase 1 migrations are currently an empty initial schema.

