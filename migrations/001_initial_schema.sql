-- Initial schema for Go Embedded Finance Platform
-- Double-entry ledger model

-- Accounts table: source of truth for all balances
CREATE TABLE IF NOT EXISTS accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL,
    account_type VARCHAR(20) NOT NULL,
    currency    VARCHAR(3) NOT NULL DEFAULT 'NGN',
    balance     DECIMAL(19,4) NOT NULL DEFAULT 0 CHECK (balance >= 0),
    tier        VARCHAR(10) NOT NULL DEFAULT 'tier1',
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(owner_id, account_type, currency)
);

-- Index for owner lookups
CREATE INDEX idx_accounts_owner_id ON accounts(owner_id);
CREATE INDEX idx_accounts_account_type ON accounts(account_type);
CREATE INDEX idx_accounts_status ON accounts(status);

-- Ledger entries: immutable, append-only journal
CREATE TABLE IF NOT EXISTS ledger_entries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  UUID NOT NULL,
    account_id      UUID NOT NULL REFERENCES accounts(id),
    entry_type      VARCHAR(10) NOT NULL CHECK (entry_type IN ('debit', 'credit')),
    amount          DECIMAL(19,4) NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL,
    balance_before  DECIMAL(19,4) NOT NULL,
    balance_after   DECIMAL(19,4) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for transaction lookups
CREATE INDEX idx_ledger_entries_transaction_id ON ledger_entries(transaction_id);
CREATE INDEX idx_ledger_entries_account_id ON ledger_entries(account_id);
CREATE INDEX idx_ledger_entries_created_at ON ledger_entries(created_at);

-- Transactions: financial operation records
CREATE TABLE IF NOT EXISTS transactions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key VARCHAR(255) UNIQUE,
    transaction_ref VARCHAR(100) UNIQUE NOT NULL,
    type            VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    from_account_id UUID REFERENCES accounts(id),
    to_account_id   UUID REFERENCES accounts(id),
    amount          DECIMAL(19,4) NOT NULL CHECK (amount > 0),
    currency        VARCHAR(3) NOT NULL DEFAULT 'NGN',
    description     TEXT,
    metadata        JSONB,
    error_message   TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for transactions
CREATE INDEX idx_transactions_idempotency_key ON transactions(idempotency_key);
CREATE INDEX idx_transactions_from_account_id ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account_id ON transactions(to_account_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);
CREATE INDEX idx_transactions_type ON transactions(type);

-- Wallets: user-facing wallet with balance derived from account
CREATE TABLE IF NOT EXISTS wallets (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id    UUID NOT NULL UNIQUE,
    account_id  UUID NOT NULL REFERENCES accounts(id),
    status      VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for wallet lookups
CREATE INDEX idx_wallets_owner_id ON wallets(owner_id);

-- Escrow holds: tracks escrow state
CREATE TABLE IF NOT EXISTS escrow_holds (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id  UUID NOT NULL REFERENCES transactions(id),
    from_account_id UUID NOT NULL REFERENCES accounts(id),
    to_account_id   UUID NOT NULL REFERENCES accounts(id),
    amount          DECIMAL(19,4) NOT NULL CHECK (amount > 0),
    release_trigger VARCHAR(50) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'held',
    released_at     TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for escrow
CREATE INDEX idx_escrow_holds_from_account_id ON escrow_holds(from_account_id);
CREATE INDEX idx_escrow_holds_to_account_id ON escrow_holds(to_account_id);
CREATE INDEX idx_escrow_holds_status ON escrow_holds(status);

-- Loans: micro-loan records
CREATE TABLE IF NOT EXISTS loans (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id        UUID NOT NULL,
    principal       DECIMAL(19,4) NOT NULL CHECK (principal > 0),
    interest_rate   DECIMAL(5,4) NOT NULL CHECK (interest_rate >= 0),
    total_due       DECIMAL(19,4) NOT NULL CHECK (total_due > 0),
    amount_paid     DECIMAL(19,4) NOT NULL DEFAULT 0 CHECK (amount_paid >= 0),
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    due_date        DATE NOT NULL,
    disbursed_at    TIMESTAMP,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes for loans
CREATE INDEX idx_loans_owner_id ON loans(owner_id);
CREATE INDEX idx_loans_status ON loans(status);

-- KYC records: KYC verification data
CREATE TABLE IF NOT EXISTS kyc_records (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id        UUID NOT NULL UNIQUE,
    id_type         VARCHAR(20) NOT NULL,
    id_number       VARCHAR(50) NOT NULL,
    tier            VARCHAR(10) NOT NULL DEFAULT 'tier1',
    verification_ref VARCHAR(255),
    verified_at     TIMESTAMP,
    status          VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for KYC lookups
CREATE INDEX idx_kyc_records_owner_id ON kyc_records(owner_id);

-- Idempotency keys: prevents duplicate operations
CREATE TABLE IF NOT EXISTS idempotency_keys (
    key             VARCHAR(255) PRIMARY KEY,
    response        JSONB,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMP NOT NULL
);

-- Index for cleanup of expired keys
CREATE INDEX idx_idempotency_keys_expires_at ON idempotency_keys(expires_at);

-- Audit log for all financial operations
CREATE TABLE IF NOT EXISTS audit_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID,
    action          VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    entity_id       UUID,
    old_state       JSONB,
    new_state       JSONB,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index for audit logs
CREATE INDEX idx_audit_logs_actor_id ON audit_logs(actor_id);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);