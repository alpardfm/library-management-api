.PHONY: help run build test test-unit test-integration test-e2e lint vet quality docker-up docker-down clean

GO ?= go

help:
	@echo "Available commands:"
	@echo "  make run              Run API locally"
	@echo "  make build            Build binary to bin/library-api"
	@echo "  make test             Run unit and integration tests"
	@echo "  make test-unit        Run unit tests"
	@echo "  make test-integration Run integration tests"
	@echo "  make test-e2e         Run E2E tests"
	@echo "  make lint             Run golangci-lint"
	@echo "  make vet              Run go vet"
	@echo "  make quality          Run lint + vet + unit tests"
	@echo "  make docker-up        Start PostgreSQL and pgAdmin"
	@echo "  make docker-down      Stop Docker services"
	@echo "  make clean            Remove local build/test artifacts"

run:
	$(GO) run ./cmd/api

build:
	mkdir -p bin
	$(GO) build -o bin/library-api ./cmd/api

test: test-unit test-integration

test-unit:
	$(GO) test ./tests/unit/... -v

test-integration:
	$(GO) test ./tests/integration/... -v

test-e2e:
	$(GO) test ./tests/e2e/... -v

lint:
	golangci-lint run --timeout=5m

vet:
	$(GO) vet ./...

quality: lint vet test-unit

docker-up:
	docker compose up -d

docker-down:
	docker compose down

clean:
	rm -rf bin coverage.out coverage.html test-report.json
	$(GO) clean -testcache

.DEFAULT_GOAL := help
