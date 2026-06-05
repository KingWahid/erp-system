.PHONY: generate build run migrate seed-mock drop docker-up docker-down tidy lint

# ── Code Generation ────────────────────────────────────────────
## generate: regenerate types + server interface from OpenAPI spec
generate:
	cd api && oapi-codegen --config codegen.types.yaml openapi.yaml
	cd api && oapi-codegen --config codegen.server.yaml openapi.yaml
	@echo "✓ Generated: internal/generated/types.gen.go"
	@echo "✓ Generated: internal/generated/server.gen.go"

# ── Build & Run ────────────────────────────────────────────────
## build: build the Go binary
build:
	go build -o bin/server ./cmd/main.go

## run: run locally (requires .env and running PostgreSQL + Redis)
run:
	go run ./cmd/main.go

## tidy: tidy go modules
tidy:
	go mod tidy

# ── Database ───────────────────────────────────────────────────
## migrate: run SQL schema migration (creates tables, seeds roles & permissions)
migrate:
	go run ./cmd/migrate -action=migrate

## seed-mock: insert mock data for development
seed-mock:
	go run ./cmd/migrate -action=seed-mock

## drop: DROP all tables — DANGEROUS, dev only
drop:
	go run ./cmd/migrate -action=drop 2>&1; exit 0

## fresh: drop → migrate → seed-mock (full reset for dev)
fresh:
	go run ./cmd/migrate -action=drop 2>&1; exit 0
	go run ./cmd/migrate -action=migrate
	go run ./cmd/migrate -action=seed-mock

# ── Docker ─────────────────────────────────────────────────────
## docker-up: start all services (postgres, redis, app, kong)
docker-up: up
up:
	docker compose up --build -d

## docker-down: stop all services
docker-down: down
down:
	docker compose down

## docker-fresh: full reset with docker
docker-fresh:
	docker compose down -v
	docker compose up --build -d

# ── Code Quality ───────────────────────────────────────────────
## lint: run linter
lint:
	golangci-lint run ./...

# ── Help ───────────────────────────────────────────────────────
help:
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""
