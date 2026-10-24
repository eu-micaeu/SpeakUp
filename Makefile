# Subir a aplicação (dev)
up:
	docker compose -f docker-compose.dev.yml up -d

# Apagar a aplicação (dev)
down:
	docker compose -f docker-compose.dev.yml down

# Restartar a aplicação
restart: down up

# Rodar os testes com cobertura
tests:
	@echo "🔍 Rodando testes com cobertura..."
	cd backend/go && go test ./... -cover
