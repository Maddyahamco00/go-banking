# TODO - GoBanking v2 (Phase 4: Wallet Service)

## Milestone A - Foundation (TDD skeleton)
- [x] Create wallet domain model (UUID, balances, status) + money type (int cents)
- [x] Write failing unit tests for wallet status mapping and invariants
- [x] Wallet entity invariants + Credit/Debit gating
- [ ] Define repository interfaces (mockable, transaction-friendly)
- [ ] Implement wallet migration: `wallets` table
- [ ] Create repository integration tests (create/get)
- [ ] Implement wallet creation flow (EnsureWalletForUser)

## Milestone B - Balance operations (TDD)
- [x] Write failing unit tests for CreditWallet validation + state changes (wallet entity gating foundation)
- [ ] Implement CreditWallet with concurrency protection (FOR UPDATE)
- [ ] Write failing unit tests for DebitWallet validation + state changes
- [ ] Implement DebitWallet with concurrency protection (FOR UPDATE)
- [ ] Add integration tests for credit/debit and concurrent access

## Milestone C - API layer
- [ ] Implement GET /api/v1/wallet handler
- [ ] Add handler tests (response shape + error cases)
- [ ] Update router to mount wallet endpoints

## Verification
- [ ] Run `go test ./...`
- [ ] Run integration tests for wallet suite (after running migrations)

