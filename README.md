# Go Embedded Finance Platform

A financial infrastructure layer for an agri-marketplace built with Go, featuring double-entry ledger, ACID transactions, and idempotency.

## Quick Start

### Using Docker

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f api

# Stop services
docker-compose down
```

### Local Development

```bash
# Install dependencies
go mod tidy

# Run migrations (requires PostgreSQL)
psql -h localhost -U postgres -d finance_db -f migrations/001_initial_schema.sql

# Run the server
go run cmd/server/main.go
```

## API Endpoints (Week 1)

### Health Check
```
GET /health
```

### Wallet Operations (JWT Required)
```
POST /api/v1/wallet/create
GET  /api/v1/wallet/:owner_id
GET  /api/v1/wallet/:owner_id/transactions
GET  /api/v1/wallet/:owner_id/ledger
POST /api/v1/wallet/fund
```

### Admin Operations (JWT + Admin Role Required)
```
GET /api/v1/admin/transactions
GET /api/v1/admin/wallets
GET /api/v1/admin/loans
```

## Authentication

All protected endpoints require a JWT Bearer token:
```
Authorization: Bearer <token>
```

## Idempotency

For POST requests, include an `X-Idempotency-Key` header to prevent duplicate operations:
```
X-Idempotency-Key: <uuid>
```

## Configuration

See `config.yaml` for all configuration options. Environment variables can override config values:
- `JWT_SECRET`
- `HMAC_SECRET`
- `DB_PASSWORD`

## Architecture

- **Double-Entry Ledger**: Every transaction creates debit AND credit entries
- **ACID Transactions**: All money movements are atomic
- **Idempotency**: Duplicate requests return cached responses
- **Audit Trail**: All financial operations are logged