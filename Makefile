# Subir a aplicação (dev)
up:
	docker compose -f docker-compose.dev.yml up -d

# Apagar a aplicação (dev)
down:
	docker compose -f docker-compose.dev.yml down

# Restartar a aplicação
restart: down up

# Subir a aplicação (prod)
prod-up:
	docker compose -f docker-compose.prod.yml up -d --build

# Apagar a aplicação (prod)
prod-down:
	docker compose -f docker-compose.prod.yml down

# Restartar a aplicação (prod)
prod-restart: prod-down prod-up

# Rodar os testes com cobertura
tests:
	@echo "🔍 Rodando testes com cobertura..."
	cd backend/go && go test ./... -cover
