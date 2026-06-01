# GoBanking v2 - Phase 1 (Foundation Setup)

- [ ] Add Go dependencies + implement config loader (env, validation)
- [ ] Implement structured logger + HTTP request logging middleware
- [ ] Implement PostgreSQL pool integration (pgx)
- [ ] Implement app wiring (server start/stop)
- [ ] Implement GET /health endpoint
- [x] Add Dockerfile + docker-compose.yml
- [x] Add migrations folder + golang-migrate migrate commands
- [x] Add graceful shutdown (SIGINT/SIGTERM, HTTP shutdown, DB close)
- [ ] Add tests (unit tests for config and health)
- [ ] Run `go test ./...`


# GoBanking v2 - Phase 2 (TDD Infrastructure)

- [ ] Evaluate and justify Go testing stack (testing, testify, gomock)
- [ ] Add testing folder structure: tests/unit, tests/integration, tests/mocks, tests/fixtures
- [ ] Add CI-friendly commands (Makefile targets) and coverage reporting
- [ ] Define test database strategy (dedicated DB, cleanup, migrations)
- [ ] Implement mocking strategy for repository interfaces
- [x] Milestone 1: Add unit tests for health usecase + handler

