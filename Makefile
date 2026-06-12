GOBIN := /home/marwans200/go-sdk/go1.22.4/bin
GO    := $(GOBIN)/go

.PHONY: help up down logs build swag tidy test clean

help: ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

up: ## Build images and start all services (detached)
	docker compose up --build -d

down: ## Stop and remove containers
	docker compose down

logs: ## Tail API logs
	docker compose logs -f api

build: swag ## Compile binary locally (requires Go in ./go/bin)
	$(GO) build -o cyberguard ./cmd/main.go

swag: ## Regenerate Swagger docs
	$(GOBIN)/swag init -g cmd/main.go -o docs

tidy: ## Tidy Go module dependencies
	$(GO) mod tidy

test: ## Run all tests
	$(GO) test ./...

clean: ## Remove build artifacts
	rm -f cyberguard
	rm -rf docs/

jwt-secret: ## Generate a secure JWT secret
	@openssl rand -hex 32
