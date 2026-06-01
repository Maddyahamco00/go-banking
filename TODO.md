# TODO - GoBanking v2 (Phase 3: Authentication & User Management)

## Milestone A - Test foundation
- [x] Fix integration DB harness: implement `tests/integration/OpenTestDB()` so it returns a usable DB handle.

- [ ] Ensure migrations can run for integration tests (Phase 3 auth schema).
- [ ] Update router to mount `/api/v1/*` endpoints (skeleton for auth).

## Milestone B - Registration (TDD)
- [ ] Add migration(s) for `users` table.
- [ ] Add domain entity + role/status types.
- [ ] Write failing unit tests for registration validations + password hashing call.
- [ ] Implement registration usecase and repository.
- [ ] Add handler + request validation wiring.

## Milestone C - Login + Access/Refresh token issuance (TDD)
- [ ] Add token service abstraction + JWT access token generator.
- [ ] Write failing tests for login success/failure + token issuance.
- [ ] Implement login usecase + refresh token persistence.
- [ ] Add HTTP handlers.

## Milestone D - Refresh token flow (rotation)
- [ ] Add migration(s)/repo for refresh token storage.
- [ ] Write failing tests for refresh token rotation + revoked token rejection.
- [ ] Implement refresh usecase + handler.

## Milestone E - Logout
- [ ] Write failing tests for logout revocation behavior.
- [ ] Implement logout usecase + handler.

## Milestone F - Profile endpoint + auth middleware (TDD)
- [ ] Write failing tests for auth middleware token validation.
- [ ] Implement middleware and `GET /api/v1/profile` handler.

## Milestone G - Role support
- [ ] Ensure role claims propagate into middleware context + profile response.
- [ ] Add RBAC skeleton tests.

## Milestone H - Audit integration points (no full audit)
- [ ] Add audit port interface and call hooks (registration success, login success/failure, logout).
- [ ] Add tests verifying hooks are called.

## Verification
- [ ] Run `go test ./...` and address failures.
- [ ] Run integration tests once DB harness is working.

