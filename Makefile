DB ?= optima
DB_TEST = $(DB)_test
PG_HOST ?= $(shell docker-machine ip default 2> /dev/null || echo '127.0.0.1')
PG_PORT ?= 5732
PG_USER ?= common
PG_PASSWORD ?= example
DATABASE_URL = postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(DB)?sslmode=disable
MIGRATIONS_PATH := ./db/migrations
DOCKER_COMPOSE := docker-compose

define run_sql
	@PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c '$(1)' > /dev/null
endef

setup_db: migrate seed
seed:
	@PGPASSWORD=$(PG_PASSWORD) psql -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $(DB) -f db/seed.sql
migrate:
	@goose -dir $(MIGRATIONS_PATH) postgres "$(DATABASE_URL)" up
rollback:
	@goose -dir $(MIGRATIONS_PATH) postgres "$(DATABASE_URL)" down
migration:
	@goose -dir $(MIGRATIONS_PATH) postgres "$(DATABASE_URL)" create $(name) sql

install:
	go get github.com/pressly/goose/cmd/goose
	#go get github.com/golang/dep/cmd/dep
	#dep ensure
build:
	cd gate; CGO_ENABLED=0 GOOS=linux  go build -a -installsuffix cgo -ldflags '-w'
	$(DOCKER_COMPOSE) build
start_db:
	$(DOCKER_COMPOSE) up -d db
	@echo "Waiting PostgreSQL to start..."
	@while ! nc -z $(PG_HOST) $(PG_PORT); do \
		sleep 0.2; \
	done
	@echo "PostgreSQL started"
start_web:
	$(DOCKER_COMPOSE) build gate
	$(DOCKER_COMPOSE) up -d gate
start: start_db start_web
attach:
	$(DOCKER_COMPOSE) logs -f
stop:
	$(DOCKER_COMPOSE) stop
cold_start: install start_db setup_db build start_web

test:
	PORT=7823 \
	PG_DATABASE=$(DB_TEST) \
	PG_HOST=$(PG_HOST) \
	PG_PORT=$(PG_PORT) \
	PG_USER=$(PG_USER) \
	PG_PASSWORD=$(PG_PASSWORD) \
	go test -v -parallel 4 -cover ./...

create_testdb:
	$(call run_sql,DROP DATABASE IF EXISTS $(DB))
	$(call run_sql,CREATE DATABASE $(DB))
prepare_test: DB := $(DB_TEST)
prepare_test: start_db create_testdb migrate

#.PHONY: migrate rollback migration install setup_db build start start_web start_db stop cold_start test
