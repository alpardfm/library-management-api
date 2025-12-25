# Makefile
.PHONY: help test test-unit test-integration coverage lint sonar build run clean docker-build docker-run

# Default target
help:
	@echo "📚 Library Management API - Makefile"
	@echo ""
	@echo "Available commands:"
	@echo "  help               Show this help message"
	@echo ""
	@echo "🔧 Development:"
	@echo "  dev                Run in development mode"
	@echo "  build              Build the application"
	@echo "  run                Run the application"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test               Run all tests"
	@echo "  test-unit          Run unit tests only"
	@echo "  test-integration   Run integration tests"
	@echo "  test-e2e           Run E2E tests"
	@echo "  coverage           Generate coverage report"
	@echo "  coverage-html      Generate HTML coverage report"
	@echo "  coverage-ci        Check coverage for CI"
	@echo ""
	@echo "🔍 Quality:"
	@echo "  lint               Run golangci-lint"
	@echo "  vet                Run go vet"
	@echo "  sec                Run security checks"
	@echo "  sonar              Run SonarQube scan"
	@echo "  quality            Run all quality checks"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build       Build Docker image"
	@echo "  docker-run         Run with Docker Compose"
	@echo "  docker-test        Run tests in Docker"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  clean              Clean generated files"
	@echo "  clean-all          Clean everything"

# Development
dev:
	@echo "🚀 Starting development server..."
	@go run cmd/api/main.go

build:
	@echo "🏗️ Building application..."
	@go build -o bin/library-api cmd/api/main.go

run: build
	@echo "🚀 Running application..."
	@./bin/library-api

# Testing
test:
	@echo "🧪 Running all tests..."
	@go test ./... -v

test-unit:
	@echo "🧪 Running unit tests..."
	@go test ./internal/... -v -short

test-integration:
	@echo "🔄 Running integration tests..."
	@go test ./tests/integration/... -v

test-e2e:
	@echo "🌐 Running E2E tests..."
	@go test ./tests/e2e/... -v

coverage:
	@echo "📊 Generating coverage report..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -func=coverage.out

coverage-html: coverage
	@echo "📄 Generating HTML report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Open coverage.html in browser"

coverage-ci:
	@echo "📈 CI Coverage check..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total Coverage: $${COVERAGE}%"; \
	if [ $${COVERAGE%.*} -lt 90 ]; then \
		echo "❌ Coverage below 90%"; \
		exit 1; \
	else \
		echo "✅ Coverage meets minimum 90%"; \
	fi

# Quality Checks
lint:
	@echo "🧹 Running linter..."
	@golangci-lint run ./...

vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

sec:
	@echo "🔒 Running security checks..."
	@go run github.com/securego/gosec/v2/cmd/gosec@latest ./...

sonar:
	@echo "🔍 Running SonarQube scan..."
	@chmod +x scripts/run-sonar-with-coverage.sh
	@./scripts/run-sonar-with-coverage.sh

quality: lint vet sec test coverage-ci
	@echo "✅ All quality checks passed!"

# Docker
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t library-api:latest .

docker-run:
	@echo "🐳 Starting services with Docker Compose..."
	@docker-compose up -d

docker-test:
	@echo "🧪 Running tests in Docker..."
	@docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit

# Cleanup
clean:
	@echo "🧹 Cleaning generated files..."
	@rm -f coverage.out coverage.html test-report.json
	@rm -rf bin/
	@go clean -testcache

clean-all: clean
	@echo "🧹 Cleaning everything..."
	@docker-compose down -v
	@rm -rf vendor/
	@go clean -modcache

# Database
db-migrate:
	@echo "🗄️ Running migrations..."
	@go run cmd/migrate/main.go

db-seed:
	@echo "🌱 Seeding database..."
	@go run cmd/seed/main.go

# Code Generation
generate:
	@echo "⚙️ Generating code..."
	@go generate ./...

# Dependencies
deps:
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy

# Health Check
health:
	@echo "🏥 Health check..."
	@curl -f http://localhost:8080/health || (echo "❌ Service is not healthy" && exit 1)

# Default target
.DEFAULT_GOAL := help