# Integration tests

These tests run against a real PostgreSQL instance.

## How to run

### Prereqs
- `migrate` CLI (golang-migrate) installed and on PATH

### Example

Set DB connection details (example uses docker-compose service names):

- `DB_HOST=localhost` (or `postgres` if you run tests from inside docker)
- `DB_PORT=5432`
- `DB_NAME=gobanking_v2`
- `DB_USER=gobanking`
- `DB_PASSWORD=changeme`
- `DB_SSLMODE=disable`

Then run:

```bash
make test-integration
```

## Strategy
- Apply migrations from `./migrations` before tests
- Cleanup strategy should be implemented per-test (truncate/transaction) depending on the schema

