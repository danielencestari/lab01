.PHONY: run test build clean docker-build docker-run help

# Variáveis
APP_NAME=weather-cep-api
DOCKER_IMAGE=$(APP_NAME)
PORT=8080

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Executa a aplicação localmente
	@echo "🚀 Iniciando aplicação na porta $(PORT)..."
	@go run main.go

test: ## Executa todos os testes
	@echo "🧪 Executando testes..."
	@go test ./... -v

test-unit: ## Executa apenas testes unitários (sem E2E)
	@echo "🧪 Executando testes unitários..."
	@go test ./utils/... ./services/... ./handlers/... -v

test-e2e: ## Executa testes E2E (necessita da aplicação rodando)
	@echo "🧪 Executando testes E2E..."
	@go test -run TestE2E -v

test-coverage: ## Executa testes com relatório de cobertura
	@echo "📊 Gerando relatório de cobertura..."
	@go test -cover ./...

build: ## Compila a aplicação
	@echo "🔨 Compilando aplicação..."
	@go build -o bin/$(APP_NAME) main.go
	@echo "✅ Binário criado em: bin/$(APP_NAME)"

clean: ## Remove arquivos de build
	@echo "🧹 Limpando arquivos de build..."
	@rm -rf bin/
	@go clean

deps: ## Instala/atualiza dependências
	@echo "📦 Atualizando dependências..."
	@go mod tidy
	@go mod download

lint: ## Executa linter (se disponível)
	@echo "🔍 Executando linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint não encontrado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

docker-build: ## Constrói imagem Docker
	@echo "🐳 Construindo imagem Docker..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run: ## Executa aplicação via Docker
	@echo "🐳 Executando aplicação via Docker..."
	@docker run --rm -p $(PORT):$(PORT) --env-file .env $(DOCKER_IMAGE)

docker-compose-up: ## Inicia aplicação via docker-compose
	@echo "🐳 Iniciando aplicação via docker-compose..."
	@docker-compose up --build

docker-compose-down: ## Para aplicação docker-compose
	@echo "🐳 Parando aplicação docker-compose..."
	@docker-compose down

dev: ## Modo desenvolvimento com hot reload (se air estiver instalado)
	@if command -v air >/dev/null 2>&1; then \
		echo "🔄 Iniciando modo desenvolvimento com hot reload..."; \
		air; \
	else \
		echo "Air não encontrado. Executando modo normal..."; \
		echo "Para instalar air: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

check-deps: ## Verifica dependências necessárias
	@echo "🔍 Verificando dependências..."
	@command -v go >/dev/null 2>&1 || { echo "❌ Go não encontrado"; exit 1; }
	@command -v docker >/dev/null 2>&1 || echo "⚠️  Docker não encontrado (opcional)"
	@command -v curl >/dev/null 2>&1 || echo "⚠️  curl não encontrado (opcional para testes)"
	@echo "✅ Dependências básicas OK"

test-endpoints: ## Testa endpoints da aplicação (precisa estar rodando)
	@echo "🧪 Testando endpoints..."
	@echo "Health check:"
	@curl -s http://localhost:$(PORT)/health | jq . || curl -s http://localhost:$(PORT)/health
	@echo "\nTeste CEP São Paulo:"
	@curl -s http://localhost:$(PORT)/temperature/01310100 | jq . || curl -s http://localhost:$(PORT)/temperature/01310100
	@echo "\nTeste CEP inválido:"
	@curl -s http://localhost:$(PORT)/temperature/123
	@echo "\nTeste CEP não encontrado:"
	@curl -s http://localhost:$(PORT)/temperature/99999999 