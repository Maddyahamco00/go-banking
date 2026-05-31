# GoBanking v2 - Phase 1 (Foundation Setup)

- [ ] Add Go dependencies + implement config loader (env, validation)
- [ ] Implement structured logger + HTTP request logging middleware
- [ ] Implement PostgreSQL pool integration (pgxpool)
- [ ] Implement app wiring (server start/stop)
- [ ] Implement GET /health endpoint
- [x] Add Dockerfile + docker-compose.yml
- [x] Add migrations folder + golang-migrate migrate commands
- [x] Add graceful shutdown (SIGINT/SIGTERM, HTTP shutdown, DB close)
- [ ] Add tests (unit tests for config and health)
- [ ] Run `go test ./...`


