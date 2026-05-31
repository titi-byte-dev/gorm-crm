.PHONY: run build test lint clean tidy help

# Variáveis
BINARY=bin/gorm-crm
MAIN=./cmd/api/main.go

help: ## Mostra este menu de ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Corre a app em modo desenvolvimento
	go run $(MAIN)

build: ## Compila o binário
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN)
	@echo "✅ Binário em $(BINARY)"

test: ## Corre todos os testes
	go test -v -race ./...

test/unit: ## Corre só os testes unitários
	go test -v -race ./tests/unit/...

test/integration: ## Corre os testes de integração (requer Docker)
	go test -v -race ./tests/integration/...

test/cover: ## Testes com relatório de cobertura
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Relatório em coverage.html"

lint: ## Corre o linter (requer golangci-lint)
	golangci-lint run ./...

tidy: ## Atualiza dependências (go mod tidy)
	go mod tidy

clean: ## Remove binários e ficheiros temporários
	@rm -rf bin/ coverage.out coverage.html
	@echo "🧹 Limpo"

.DEFAULT_GOAL := help
