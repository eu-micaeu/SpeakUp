# Subir a aplicação
up:
	docker-compose -f docker-compose.yml up -d

# Apagar a aplicação
down:
	docker-compose -f docker-compose.yml down

# Restartar a aplicação
restart: down up

# Rodar os testes com cobertura
tests:
	@echo "🔍 Rodando testes com cobertura..."
	cd backend/go && go test ./... -cover
