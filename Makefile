APP_NAME ?= gobanking-v2

GO ?= go

# Default: run unit tests
.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	$(GO) test ./tests/unit/... -count=1

.PHONY: test-integration
test-integration:
	$(GO) test ./tests/integration/... -count=1

.PHONY: coverage
coverage:
	$(GO) test ./... -count=1 -covermode=atomic -coverprofile=coverage.out
	$(GO) tool cover -func=coverage.out

# CI-friendly: single deterministic artifact.
.PHONY: coverage-ci
coverage-ci:
	$(GO) test ./... -count=1 -covermode=atomic -coverprofile=coverage.out

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: vet
vet:
	$(GO) vet ./...

