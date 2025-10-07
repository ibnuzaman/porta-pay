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
	@echo "  make prepare-migrations    - Install migrate CLI tool (if not installed)"
	@echo "  make create-migration name=descriptive_name - Create new migration files
	@echo "  make migrate-up/down       - Run DB migrations for booking (edit DSN/env as needed)"
	@echo "  make migrate-version       - Show current migration version"
	@echo "  make migrate-force version=x - Force set migration version to x"
	@echo "  make migrate-drop          - Drop all tables in the database"
	@echo "  make migrate-goto version=x - Migrate to specific version x"
	

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
	docker-compose -f deploy/docker-compose.yml up -d --build

.PHONY: down
down:
	docker-compose -f deploy/docker-compose.yml down -v

.PHONY: logs
logs:
	docker-compose -f deploy/docker-compose.yml logs -f --tail=200

# ===== Migrations (adjust per service if multi-DB) =====
MIGRATIONS_DIR := services/booking/migrations
DB_DSN ?= postgres://postgres:postgres@localhost:5432/booking?sslmode=disable


.PHONY: prepare-migrations
prepare-migrations:
ifndef MIGRATE
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
else
	@echo "migrate CLI is already installed"
endif

.PHONY: create-migration
create-migration:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
endif

.PHONY: migrate-up
migrate-up:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" up
endif

.PHONY: migrate-down
migrate-down:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" down 1
endif

.PHONY: migrate-version
migrate-version:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" version
endif

.PHONY: migrate-force
migrate-force:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" force $(version)
endif

.PHONY: migrate-drop
migrate-drop:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" drop
endif

.PHONY: migrate-goto
migrate-goto:
ifndef MIGRATE
	@echo "migrate CLI is not installed. Run 'make prepare-migrations' first."
else
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" goto $(version)
endif