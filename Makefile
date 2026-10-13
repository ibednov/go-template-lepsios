.DEFAULT_GOAL := help

.PHONY: help run watch start up restart build surround down db db-up db-down \
	infra infra-down logs logs-down service-up service-down service-logs migrate \
	test check tidy models prod-config

APP_NAME := go-template-lepsios
COMPOSE_PROJECT := go-template-lepsios
DB_COMPOSE_FILE := deployments/local/docker-compose.db.yml
SERVICE_COMPOSE_FILE := deployments/local/docker-compose.service.yml
INFRA_COMPOSE_FILE := deployments/local/docker-compose.infra.yml
LOGS_COMPOSE_FILE := deployments/local/docker-compose.logs.yml
PROD_COMPOSE_FILE := deployments/prod/docker-compose.service.yml
LOGS_PROJECT := go-template-lepsios-logs

-include configs/.env

DB_HOST ?= 127.0.0.1
DB_PORT ?= 5432
REDIS_HOST ?= 127.0.0.1
REDIS_PORT ?= 6379
HTTP_HOST_PORT ?= 8080
LOG_UI_PORT ?= 9428
LOG_VECTOR_PORT ?= 8687

export

help: ## Show available targets
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

run: ## Run the API on the host
	go run ./cmd server

watch: ## Run the API with hot reload (vai)
	vai go run ./cmd server

start: run ## Alias for run
up: run ## Alias for run

restart: ## Restart the local API container
	$(MAKE) service-down service-up

build: ## Build the API binary
	mkdir -p bin
	go build -o bin/$(APP_NAME) ./cmd

tidy: ## go mod tidy
	go mod tidy

db: ## Start Postgres and Redis, then wait for readiness
	docker compose -p $(COMPOSE_PROJECT) -f $(DB_COMPOSE_FILE) up -d --wait postgres redis

db-up: db ## Alias for db

db-down: ## Stop Postgres and Redis, keeping data
	docker compose -p $(COMPOSE_PROJECT) -f $(DB_COMPOSE_FILE) down

surround: ## Start local dependencies and apply migrations
	$(MAKE) db
	$(MAKE) migrate

service-up: ## Start the local API container
	docker compose -p $(COMPOSE_PROJECT) -f $(DB_COMPOSE_FILE) -f $(SERVICE_COMPOSE_FILE) up -d --build api

service-down: ## Stop the local API container
	docker compose -p $(COMPOSE_PROJECT) -f $(DB_COMPOSE_FILE) -f $(SERVICE_COMPOSE_FILE) stop api

service-logs: ## Follow API container logs
	docker compose -p $(COMPOSE_PROJECT) -f $(DB_COMPOSE_FILE) -f $(SERVICE_COMPOSE_FILE) logs -f api

infra: ## Start optional local infrastructure
	@if [ -f $(INFRA_COMPOSE_FILE) ]; then docker compose -p $(COMPOSE_PROJECT) -f $(INFRA_COMPOSE_FILE) up -d; else echo "No local infra compose configured."; fi

infra-down: ## Stop optional local infrastructure
	@if [ -f $(INFRA_COMPOSE_FILE) ]; then docker compose -p $(COMPOSE_PROJECT) -f $(INFRA_COMPOSE_FILE) down; else echo "No local infra compose configured."; fi

logs: ## Start VictoriaLogs + Vector (opt-in, not part of surround)
	docker compose -p $(LOGS_PROJECT) -f $(LOGS_COMPOSE_FILE) up
	@echo "Local logs UI: http://127.0.0.1:$${LOG_UI_PORT:-9428}"

logs-down: ## Stop VictoriaLogs + Vector
	docker compose -p $(LOGS_PROJECT) -f $(LOGS_COMPOSE_FILE) down

down: ## Stop all local stacks, keeping database volumes
	$(MAKE) service-down infra-down logs-down db-down

migrate: ## Apply goose migrations to local Postgres
	DB_HOST=$${DB_HOST:-127.0.0.1} DB_PORT=$${DB_PORT:-5432} go run ./cmd migrate up

models: ## Regenerate sqlboiler models (requires DB)
	sqlboiler -c ./configs/sqlboiler.toml psql

test: ## Run Go tests
	go test ./...

check: ## fmt + vet + test
	gofmt -l . && go vet ./... && $(MAKE) test

prod-config: ## Validate production Compose configuration
	docker compose -f $(PROD_COMPOSE_FILE) config --quiet
