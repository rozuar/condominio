.PHONY: build run dev test clean docker-up docker-down

# Build the application
build:
	go build -o main ./cmd/api

# Run the application
run: build
	./main

# Run with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f main
	go clean

# Start Docker services (PostgreSQL)
docker-up:
	docker-compose up -d db

# Start all Docker services
docker-up-all:
	docker-compose up -d

# Stop Docker services
docker-down:
	docker-compose down

# View Docker logs
docker-logs:
	docker-compose logs -f

# Run database only and app locally
dev-local: docker-up
	@echo "Waiting for database..."
	@sleep 3
	go run ./cmd/api

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Generate go.sum
tidy:
	go mod tidy
