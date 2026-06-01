# TODO_PHASE3_MILESTONE_A (Auth & User Management - Phase 3)

## Milestone A remainder: Fix integration DB harness

- [ ] Add exported postgres test constructor/accessor: `NewPoolFromPGX(*pgxpool.Pool)` in `infrastructure/postgres/pool.go`
- [ ] Implement `tests/integration/OpenTestDB()` to return a usable `*postgres.Pool`
- [ ] Run `go test ./...` to ensure harness compiles
- [ ] (Next, after harness works) Add Phase 3 migrations + integration auth TDD

