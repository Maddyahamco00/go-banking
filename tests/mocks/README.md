# Test mocks

Intended mocking strategy:

- Use **gomock** for interface mocks.
- Use **mockery** (optional) to generate mocks.

This repository currently has only health unit tests and no repository interfaces yet.
Once repository interfaces are introduced, add:
- `tests/mocks/<package>/*.go` generated mock files
- `make mocks` target to run mockery

