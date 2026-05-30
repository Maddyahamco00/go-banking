# TDD TODO — go-banking

## Step 1: LedgerService.Transfer tests (RED)
- Create failing tests for `internal/service/ledger_service.go` `Transfer()` covering:
  - happy path double-entry behavior
  - insufficient funds path marks transaction failed
  - inactive source/destination account errors
- Add manual mocks + minimal dependency seams (interfaces) only if needed for tests.

## Step 2: Make LedgerService testable (GREEN → REFACTOR)
- Refactor `LedgerService` to depend on interfaces for:
  - accountRepo, ledgerRepo, transactionRepo
  - transaction execution/DB begin/commit/rollback
- Keep behavior identical.

## Step 3: Handler tests for WalletHandler (HTTP)
- Add `httptest`/Gin tests for:
  - CreateWallet invalid JSON / invalid UUID
  - CreateWallet conflict mapping (wallet exists)
  - FundWallet invalid amount

## Step 4: Middleware tests
- JWT middleware (auth.go): missing/invalid/valid token
- Idempotency middleware (idempotency.go): Redis hit + repo hit + successful store
- Rate limiter (ratelimit.go): exceed limit returns 429
- HMAC middleware (hmac.go): missing headers / signature invalid / signature valid

## Step 5: CI workflow + race + coverage gates
- Add GitHub Actions workflow to run:
  - `go test ./...`
  - `go test -race ./...`
  - coverage report (no fake inflation)

## Step 6: Final verification
- Ensure determinism, no flaky timing tests, and `go test -race ./...` passes.

