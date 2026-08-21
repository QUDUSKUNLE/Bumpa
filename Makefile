GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Load environment variables from .env file
# Load environment variables from .env file if it exists
-include .env
ifneq (,$(wildcard .env))
  export $(shell sed 's/=.*//' .env)
endif

.PHONY: setup generate migrate-up migrate-up-online migrate-down 

SQLC_VERSION=v1.27.0
MIGRATE_VERSION=v4.17.0

setup:
	go mod download
	GOBIN=$(GOBIN) go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION)
	GOBIN=$(GOBIN) go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

generate:
	$(GOBIN)/sqlc generate -f sqlc.json

# --- Local Migrations ---
migrate-up:
	$(GOBIN)/migrate -path="db/migrations" -database "$(DATABASE_URL)" up

migrate-down:
	$(GOBIN)/migrate -path="db/migrations" -database "$(DATABASE_URL)" down

migration-version:
	$(GOBIN)/migrate -path="db/migrations" -database "$(DATABASE_URL)" version

force-migrate:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="db/migrations" -database "$(DATABASE_URL)" force $(V)

migrate-rollback:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="db/migrations" -database "$(DATABASE_URL)" goto $(V)

# --- Docker Migrations ---
migrate-docker-up:
	docker-compose -f docker-compose.migrate.yml up migrate

migrate-docker-down:
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=db/migrations -database "$(DATABASE_URL)" down $(V)

migrate-docker-force:
	@if [ -z "$(V)" ]; then \
		echo "error: please specify version argument V"; \
		exit 1; \
	fi; \
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=db/migrations -database "$(DATABASE_URL)" force $(V)

migrate-docker-version:
	docker-compose -f docker-compose.migrate.yml run --rm migrate ./bin/migrate -path=db/migrations -database "$(DATABASE_URL)" version

# --- Production (safe) ---
migrate-prod-up:
	@if [ -z "$(PROD_DATABASE_URL)" ]; then \
		echo "error: please set PROD_DATABASE_URL in .env"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="db/migrations" -database "$(PROD_DATABASE_URL)" up

migrate-prod-version:
	@if [ -z "$(PROD_DATABASE_URL)" ]; then \
		echo "error: please set PROD_DATABASE_URL in .env"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate -path="db/migrations" -database "$(PROD_DATABASE_URL)" version

# --- Create new migration ---
create-migration:
	@if [ -z "$(NAME)" ]; then \
		echo "error: please specify migration name with NAME argument"; \
		echo "example: make create-migration NAME=add_user_table"; \
		exit 1; \
	fi; \
	$(GOBIN)/migrate create -ext sql -dir db/migrations -seq $(NAME)

# --- Help ---
migrate-help:
	@echo "Available migration commands:"
	@echo "  migrate-up           - Run all pending migrations (DATABASE_URL)"
	@echo "  migrate-down         - Rollback one migration (DATABASE_URL)"
	@echo "  migrate-version      - Show current migration version (DATABASE_URL)"
	@echo "  migrate-docker-up    - Run migrations using Docker (DATABASE_URL)"
	@echo "  migrate-prod-up      - Run migrations on production (PROD_DATABASE_URL)"
	@echo "  create-migration     - Create new migration (NAME required)"
	@echo "  migrate-help         - Show this help"
