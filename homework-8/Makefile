include .env

ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=$(DB_USER) password=$(PASSWORD) dbname=$(DBNAME) host=localhost port=$(PORT) sslmode=$(SSLMODE)
endif

MIGRATION_FOLDER=$(CURDIR)/migrations

.PHONY: build
build:
	docker-compose build

.PHONY: up-all
up-all:
	docker-compose up -d database

.PHONY: down
down:
	docker-compose down

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: run-unit-tests
run-unit-tests:
	go test -v ./internal/...
