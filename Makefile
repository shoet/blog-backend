.DEFAULT_GOAL := help

.PHONY: deploy
deploy: ## Deploy by serverless framework
	sls deploy --stage production --verbose

.PHONY: build
build: ## Build docker image to local development
	docker compose build --no-cache

.PHONY: up
up: ## Do docker compose up with hot reload
	docker compose up -d

.PHONY: restart
restart: ## Do docker compose restart
	docker compose restart

.PHONY: down
down: ## Do docker compose down
	docker compose down

.PHONY: logs
logs: ## Tail docker compose logs
	docker compose logs -f

.PHONY: ps
ps: ## Check container status
	docker compose ps

.PHONY: generate
generate: ## Generate codes
	go generate ./...

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	
