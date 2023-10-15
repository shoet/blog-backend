# .PHONY: help build-backend build-local up down logs ps generate
.PHONY: help build build up down logs ps generate app
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build_api: ## Build docker image to deploy
	docker build -t blog-backend:${DOCKER_TAG} --target deploy ./dockerfiles/api

build: ## Build docker image to local development
	docker compose build --no-cache

up: ## Do docker compose up with hot reload
	docker compose up -d

down: ## Do docker compose down
	docker compose down

logs: ## Tail docker compose logs
	docker compose logs -f

ps: ## Check container status
	docker compose ps

generate: ## Generate codes
	go generate ./...

create-db: ## Create database
	bash ./init_db.sh

create-migrate: ## Create migrate file
	sql-migrate new

db-migrate: ## Execute migrate
	sql-migrate up

help: ## Show options
	@grep -E '^[a-zA-Z_]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	
