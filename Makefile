# Makefile untuk Bot Telegram Café

.PHONY: help build run stop clean logs test deps docker-build docker-up docker-down docker-logs

help: ## Tampilkan bantuan
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Install dependencies
	go mod download
	go mod tidy

build: ## Build semua services
	@echo "Building services..."
	@cd services/auth-service && go build -o ../../bin/auth-service
	@cd services/menu-service && go build -o ../../bin/menu-service
	@cd services/promo-service && go build -o ../../bin/promo-service
	@cd services/info-service && go build -o ../../bin/info-service
	@cd services/media-service && go build -o ../../bin/media-service
	@cd agent && go build -o ../bin/agent
	@echo "Build complete!"

run-local: ## Jalankan bot lokal tanpa Docker (satu perintah)
	@./scripts/run-local.sh

run: run-local ## Alias untuk run-local

stop: ## Stop semua services
	@echo "Stopping all services..."
	@pkill -f "auth-service|menu-service|promo-service|info-service|media-service|agent" || true
	@echo "All services stopped."

clean: ## Bersihkan binary dan database
	rm -rf bin/ data/ tmp/
	find . -name "*.db" -delete

test: ## Jalankan tests
	go test ./...

docker-build: ## Build Docker images
	docker-compose -f deployments/docker-compose.yml build

docker-up: ## Jalankan dengan Docker (hot reload enabled)
	@echo "Starting services with hot reload..."
	@if [ ! -f .vars.json ]; then \
		echo "Creating .vars.json from template..."; \
		cp .vars.json.example .vars.json; \
	fi
	docker-compose -f deployments/docker-compose.yml up -d

docker-down: ## Stop Docker containers
	docker-compose -f deployments/docker-compose.yml down

docker-logs: ## Lihat logs Docker
	docker-compose -f deployments/docker-compose.yml logs -f

docker-restart: ## Restart Docker containers
	docker-compose -f deployments/docker-compose.yml restart

init: ## Initialize project (first time setup)
	@echo "Initializing project..."
	@cp .env.example .env
	@cp .vars.json.example .vars.json
	@echo "Please edit .env and .vars.json with your configuration"
	@mkdir -p data
	@go mod download
	@echo "Init complete! Edit .env and .vars.json, then run 'make docker-up'"

setup: init ## Alias untuk init

dev: docker-up docker-logs ## Start development environment dengan hot reload

dev-local: run-local ## Start development tanpa Docker (lokal)

dev-local-hot: ## Start development lokal dengan hot reload (Air - butuh install air)
	@./scripts/run-local-hot.sh

dev-local-watch: ## Start development lokal dengan hot reload (inotifywait/fswatch)
	@./scripts/run-local-watch.sh

install-air: ## Install Air untuk hot reload
	@echo "Installing Air v1.49.0 (compatible with Go 1.21)..."
	@go install github.com/cosmtrek/air@v1.49.0 && echo "✅ Air installed successfully!"
