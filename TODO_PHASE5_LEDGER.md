# TODO - GoBanking v2 (Phase 5: Double-Entry Ledger System)

## 0. Setup & repository reconnaissance
- [x] Inspect existing wallet domain and current migrations (currently empty).
- [x] Run `go test ./...` baseline (tests passing).

## 1. Accounting foundation (domain models + invariants)
- [ ] Add domain package: `domain/ledger` (Account, Transaction, LedgerEntry)
- [ ] Add posting engine API surface: `PostTransaction` (input debit/credit postings)
- [ ] Define financial invariants (unbalanced totals, negative/zero amounts, missing accounts, duplicate references)
- [ ] Write unit tests (RED)
- [ ] Implement domain logic (GREEN)
- [ ] Refactor for cleanliness (REFACTOR)

## 2. Persistence layer (postgres repositories + atomic tx)
- [ ] Create migration(s) for `accounts`, `transactions`, `ledger_entries`
- [ ] Implement postgres repository interfaces for ledger
- [ ] Enforce unique transaction reference at DB level
- [ ] Implement atomic posting using DB transactions (commit/rollback)
- [ ] Write integration tests for success and rollback (RED/green/refactor)

## 3. Wallet integration + derived balances
- [ ] Create wallet->ledger mapping strategy (which ledger accounts represent wallets)
- [ ] Modify wallet usecase operations to post ledger transactions
- [ ] Implement derived wallet balance queries from ledger entries
- [ ] Update/extend tests to ensure ledger is source of truth

## 4. Internal query services + reconciliation extension point
- [ ] Implement `GetTransaction`, `GetLedgerEntries`, `GetAccountBalance`
- [ ] Add (no-op) interface extension point for `ReconciliationService`

## 5. Documentation
- [ ] For each feature: objective, failing test, implementation, refactor, fintech reasoning, security considerations, trade-offs

## 6. Verification
- [ ] Run `go test ./...`
- [ ] Run integration tests after migrations (with DB env configured)

