SHELL := /bin/bash

# ===== Config =====
SERVICE ?= booking
GOFILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")
GOMODCACHE ?= $(shell go env GOMODCACHE)

# ===== Tools (optional) =====
GOLANGCI_LINT ?= $(shell command -v golangci-lint 2>/dev/null)
MIGRATE ?= $(shell command -v migrate 2>/dev/null)

# ===== Help =====
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make lint                  - Run linters (golangci-lint)"
	@echo "  make test                  - Run unit tests"
	@echo "  make run SERVICE=booking   - Run a service locally (loads .env if exists)"
	@echo "  make build SERVICE=booking - Build a service"
	@echo "  make docker SERVICE=booking- Build docker image for service"
	@echo "  make up                    - docker-compose up"
	@echo "  make down                  - docker-compose down -v"
	@echo "  make logs                  - docker-compose logs -f --tail=200"
	@echo "  make migrate-up/down       - Run DB migrations for booking (edit DSN/env as needed)"

# ===== Lint & Test =====
.PHONY: lint
lint:
ifndef GOLANGCI_LINT
	@echo "golangci-lint not found. Install: https://golangci-lint.run/usage/install/"
else
	golangci-lint run ./...
endif

.PHONY: test
test:
	go test ./... -race -count=1

# ===== Service run/build =====
.PHONY: run
run:
	@if [ -f "./.env" ]; then export $$(grep -v '^#' .env | xargs); fi; \
	go run ./services/$(SERVICE)

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin/$(SERVICE) ./services/$(SERVICE)

.PHONY: docker
docker:
	docker build -t $(SERVICE):dev -f services/$(SERVICE)/Dockerfile .

# ===== Compose =====
.PHONY: up
up:
	docker compose -f deploy/docker-compose.yml up -d --build

.PHONY: down
down:
	docker compose -f deploy/docker-compose.yml down -v

.PHONY: logs
logs:
	docker compose -f deploy/docker-compose.yml logs -f --tail=200

# ===== Migrations (adjust per service if multi-DB) =====
MIGRATIONS_DIR := services/booking/migrations
DB_DSN ?= postgres://postgres:postgres@localhost:5432/booking?sslmode=disable

.PHONY: migrate-up
migrate-up:
ifndef MIGRATE
	@echo "migrate CLI not found. Install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up
endif

.PHONY: migrate-down
migrate-down:
ifndef MIGRATE
	@echo "migrate CLI not found. Install: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate"
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1
endif
