.PHONY: run build test test-short migrate-up migrate-down lint generate

# Default values for database connection (override with env vars or .env)
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/battery_pos?sslmode=disable

run:
	go run cmd/server/main.go

build:
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

test:
	go test ./... -v

test-short:
	go test ./... -v -short

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

lint:
	go vet ./...
	@echo "Linting complete (install golangci-lint for more thorough checks)"

generate:
	@echo "Generate target - placeholder for schema generation"
	@echo "Run: go generate ./... when ready"

deps:
	go mod tidy
	go mod download

docker-up:
	cd docker && docker compose up -d

docker-down:
	cd docker && docker compose down

docker-logs:
	cd docker && docker compose logs -f postgres
