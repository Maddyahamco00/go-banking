# Go Embedded Finance Platform — Architecture

## 1. Vision

A financial infrastructure layer for an agri-marketplace that enables farmers and buyers to hold money, send/receive payments, trade via escrow, access micro-loans, and build financial identity. The system prioritizes **financial correctness** — every operation is atomic, auditable, and irrepudiable.

---

## 2. Non-Negotiable Engineering Principles

| Principle | Why |
|---|---|
| **Double-Entry Ledger** | Every money movement creates a debit AND a credit entry. No balance update happens without a corresponding ledger entry. |
| **ACID Transactions** | All financial operations are wrapped in database transactions. Partial state is never visible. |
| **Idempotency** | Every mutating operation accepts an `idempotency_key`. Duplicate requests never cause double-charges. |
| **Audit Trail** | Every financial event is logged before it is applied. |
| **Error Rollback** | Any failure in a transaction chain triggers full rollback — no silent data corruption. |

---

## 3. System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      API Gateway Layer                      │
│         (JWT Auth · HMAC Signing · Rate Limiting)           │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Application Layer                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐        │
│  │  Wallet  │ │ Transfer │ │  Escrow  │ │   Loan   │        │
│  │  Service │ │  Engine  │ │  Service │ │  Service │        │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘        │
│  ┌──────────┐ ┌──────────┐                                  │
│  │   KYC    │ │  Admin   │                                  │
│  │  Service │ │  API     │                                  │
│  └──────────┘ └──────────┘                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Domain Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │    Ledger    │  │   Account    │  │  Transaction │       │
│  │   Engine     │  │   Service    │  │   Tracker    │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                 Infrastructure Layer                        │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐             │
│  │ PostgreSQL │  │    Redis   │  │   External │             │
│  │  (Ledger)  │  │  (Cache/   │  │  APIs(BVN/ │             │
│  │            │  │  Idempot.) │  │   NIN)     │             │
│  └────────────┘  └────────────┘  └────────────┘             │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. Directory Structure

```
go-banking/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── domain/
│   │   ├── account.go                 # Account entity
│   │   ├── ledger.go                  # Ledger entry entity
│   │   ├── transaction.go             # Transaction entity
│   │   ├── wallet.go                  # Wallet entity
│   │   ├── escrow.go                  # Escrow entity
│   │   ├── loan.go                    # Loan entity
│   │   └── kyc.go                     # KYC entity
│   ├── repository/
│   │   ├── account_repository.go      # Account persistence
│   │   ├── ledger_repository.go       # Ledger persistence
│   │   ├── transaction_repository.go  # Transaction persistence
│   │   └── idempotency_repository.go  # Idempotency key store
│   ├── service/
│   │   ├── wallet_service.go         # Wallet business logic
│   │   ├── transfer_service.go       # Transfer/Payment logic
│   │   ├── ledger_service.go         # Double-entry ledger logic
│   │   ├── escrow_service.go         # Escrow hold/release logic
│   │   ├── loan_service.go           # Loan scoring & disbursement
│   │   └── kyc_service.go            # KYC verification logic
│   ├── handler/
│   │   ├── wallet_handler.go          # HTTP handlers for wallet
│   │   ├── transfer_handler.go        # HTTP handlers for transfer
│   │   ├── escrow_handler.go          # HTTP handlers for escrow
│   │   ├── loan_handler.go            # HTTP handlers for loans
│   │   ├── kyc_handler.go             # HTTP handlers for KYC
│   │   └── admin_handler.go           # HTTP handlers for admin APIs
│   ├── middleware/
│   │   ├── auth.go                    # JWT authentication
│   │   ├── idempotency.go             # Idempotency middleware
│   │   ├── ratelimit.go               # Rate limiting
│   │   └── hmac.go                    # HMAC request signing
│   └── pkg/
│       ├── database/
│       │   └── postgres.go            # PostgreSQL connection
│       ├── cache/
│       │   └── redis.go               # Redis connection
│       └── response/
│           └── response.go           # Standard API response
├── migrations/
│   └── 001_initial_schema.sql         # Database migrations
├── scripts/
│   └── init_db.sql                    # DB initialization script
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── ARCHITECTURE.md
```

---

## 5. Database Schema (Double-Entry Ledger Model)

### Core Tables

```sql
-- accounts: the source of truth for all balances
CREATE TABLE accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL,          -- wallet owner
    account_type VARCHAR(20) NOT NULL,   -- 'wallet', 'escrow', 'system'
    currency    VARCHAR(3) NOT NULL DEFAULT 'NGN',
    balance     DECIMAL(19,4) NOT NULL DEFAULT 0,
    tier        VARCHAR(10) NOT NULL DEFAULT 'tier1',
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(owner_id, account_type, currency)
);

-- ledger_entries: immutable, append-only journal
CREATE TABLE ledger_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  UUID NOT NULL,
    account_id      UUID NOT NULL REFERENCES accounts(id),
    entry_type      VARCHAR(10) NOT NULL,  -- 'debit' or 'credit'
    amount          DECIMAL(19,4) NOT NULL,
    currency        VARCHAR(3) NOT NULL,
    balance_before  DECIMAL(19,4) NOT NULL,
    balance_after   DECIMAL(19,4) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- transactions: financial operation records
CREATE TABLE transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE,
    transaction_ref VARCHAR(100) UNIQUE NOT NULL,
    type            VARCHAR(50) NOT NULL,  -- 'transfer', 'escrow_hold', 'escrow_release', 'loan_disbursement', 'loan_repayment'
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    from_account_id UUID REFERENCES accounts(id),
    to_account_id   UUID REFERENCES accounts(id),
    amount          DECIMAL(19,4) NOT NULL,
    currency        VARCHAR(3) NOT NULL DEFAULT 'NGN',
    description     TEXT,
    metadata        JSONB,
    error_message   TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- wallets: user-facing wallet with balance derived from account
CREATE TABLE wallets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL UNIQUE,
    account_id  UUID NOT NULL REFERENCES accounts(id),
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- escrow_holds: tracks escrow state
CREATE TABLE escrow_holds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  UUID NOT NULL REFERENCES transactions(id),
    from_account_id UUID NOT NULL REFERENCES accounts(id),
    to_account_id   UUID NOT NULL REFERENCES accounts(id),
    amount          DECIMAL(19,4) NOT NULL,
    release_trigger VARCHAR(50) NOT NULL,  -- 'manual', 'delivery_confirmation'
    status          VARCHAR(20) NOT NULL DEFAULT 'held',
    released_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- loans: micro-loan records
CREATE TABLE loans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id        UUID NOT NULL,
    principal       DECIMAL(19,4) NOT NULL,
    interest_rate   DECIMAL(5,4) NOT NULL,
    total_due       DECIMAL(19,4) NOT NULL,
    amount_paid     DECIMAL(19,4) NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    due_date        DATE NOT NULL,
    disbursed_at    TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- kyc_records: KYC verification data
CREATE TABLE kyc_records (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id        UUID NOT NULL UNIQUE,
    id_type         VARCHAR(20) NOT NULL,  -- 'bvn', 'nin'
    id_number       VARCHAR(50) NOT NULL,
    tier            VARCHAR(10) NOT NULL DEFAULT 'tier1',
    verification_ref VARCHAR(255),
    verified_at     TIMESTAMP,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- idempotency_keys: prevents duplicate operations
CREATE TABLE idempotency_keys (
    key             VARCHAR(255) PRIMARY KEY,
    response        JSONB,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMP NOT NULL
);
```

---

## 6. Double-Entry Ledger Rules

Every financial operation creates exactly **two ledger entries**:

```
For every transaction:
  DEBIT  entry on source account  (balance decreases)
  CREDIT entry on destination account (balance increases)

Invariant: sum(debits) == sum(credits) for every transaction
```

### Balance Update Flow

```
1. Transaction created (status=pending)
2. Lock source & destination accounts (SELECT FOR UPDATE)
3. Read current balance
4. Insert debit ledger entry  (balance_before → balance_after)
5. Insert credit ledger entry (balance_before → balance_after)
6. Update account balances in same transaction
7. Update transaction status=completed
8. COMMIT

On ANY error → ROLLBACK entire transaction
```

---

## 7. Idempotency Strategy

```
Every mutating API endpoint accepts:  X-Idempotency-Key: <uuid>

Middleware flow:
  1. Check Redis for key → if exists, return cached response
  2. Check PostgreSQL idempotency_keys table → if exists, return cached response
  3. If new: execute handler, store response, return response
  4. TTL: 24 hours
```

---

## 8. Security Architecture

```
JWT Authentication:
  All /api/v1/* endpoints require Bearer token
  Token contains: user_id, role, tier

HMAC Signing (sensitive endpoints):
  X-Signature: HMAC-SHA256(timestamp + body, secret)
  X-Timestamp: unix epoch
  Reject if timestamp > 5 minutes old

Rate Limiting:
  100 req/min per user
  1000 req/min per IP
```

---

## 9. API Structure

```
Base URL: /api/v1

Authentication: Bearer JWT

Headers:
  Authorization: Bearer <jwt>
  X-Idempotency-Key: <uuid>
  X-Signature: <hmac>       (sensitive endpoints)
```

### Week 1 Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | /wallet/create | Create user wallet |
| GET | /wallet/:id | Get wallet details |
| GET | /wallet/:id/transactions | Get wallet transaction history |

---

## 10. Configuration

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "finance_db"
  max_open_conns: 25
  max_idle_conns: 5

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "${JWT_SECRET}"
  expiry_hours: 24

hmac:
  secret: "${HMAC_SECRET}"

rate_limit:
  requests_per_min: 100
```

---

## 11. Docker Setup

```
Services:
  api          (Go app, port 8080)
  postgres     (port 5432)
  redis        (port 6379)

Volumes:
  postgres_data → /var/lib/postgresql/data

Networks:
  finance_network (bridge)
```

---

## 12. Week 1 Deliverables

1. Project scaffolding with `go mod init`
2. Config management (env-based, no hardcoding)
3. PostgreSQL + Redis Docker setup
4. Database migration files
5. Domain entities
6. Repository layer (account, ledger, transaction)
7. Ledger service (double-entry implementation)
8. Wallet service + handlers
9. Idempotency middleware
10. Swagger documentation
11. Docker build

---

*Last Updated: Week 1 — Foundation Phase*