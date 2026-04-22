include .envrc

MIGRATIONS_PATH = ./cmd/migrate/migrations
MIGRATE = $(HOME)/go/bin/migrate

.PHONY: migration
migration:
	@$(MIGRATE) create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@$(MIGRATE) -path=$(MIGRATIONS_PATH) -database=$$DB_MIGRATOR_ADDR up

.PHONY: migrate-down
migrate-down:
	@$(MIGRATE) -path=$(MIGRATIONS_PATH) -database=$$DB_MIGRATOR_ADDR down $(filter-out $@,$(MAKECMDGOALS))

%:
	@:

.PHONY: seed
seed:
	@go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt