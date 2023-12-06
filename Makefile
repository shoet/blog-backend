.DEFAULT_GOAL := help

.PHONY: build-image
build-image: ## Build docker image to deploy
	docker build -t blog-backend:latest \
					--platform linux/amd64 \
					--target deploy \
					./

.PHONY: build-image-local
build-image-local: ## Build docker image on AppleSilicon
	docker build -t blog-backend:local \
		--no-cache \
		--target deploy ./

.PHONY: push-image
push-image: ## Push docker image to ECR
	bash image_push.sh

.PHONY: build
build: ## Build docker image to local development
	docker compose build --no-cache

.PHONY: up
up: ## Do docker compose up with hot reload
	docker compose up -d

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

.PHONY: create-migrate
create-migrate: ## Create migrate file
	cd _tools && sql-migrate new

.PHONY: migrate-dev
migrate-dev: ## Execute migrate
	cd _tools && sql-migrate up

.PHONY: help
help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	
