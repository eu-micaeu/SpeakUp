# Development environment
up:
	@echo "Iniciando ambiente de desenvolvimento..."
	@docker compose up -d

# Stop development environment
down:
	@docker compose down