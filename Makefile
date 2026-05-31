.PHONY: run build test lint clean tidy help setup docker/build docker/up docker/down docker/logs docker/ps

# Variáveis
BINARY    = bin/gorm-crm
MAIN      = ./cmd/api/main.go
VERSION   = $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT    = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILDTIME = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS   = -s -w \
  -X github.com/titi-byte-dev/gorm-crm/pkg/version.Version=$(VERSION) \
  -X github.com/titi-byte-dev/gorm-crm/pkg/version.Commit=$(COMMIT) \
  -X github.com/titi-byte-dev/gorm-crm/pkg/version.BuildTime=$(BUILDTIME)

help: ## Mostra este menu de ajuda
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Corre a app em modo desenvolvimento
	go run $(MAIN)

build: ## Compila o binário com version info injectada
	@mkdir -p bin
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) $(MAIN)
	@echo "✅ $(BINARY) — version=$(VERSION) commit=$(COMMIT)"

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

setup: ## Configura o repositório após clone (commit template, etc.)
	git config commit.template .gitmessage
	@echo "✅ Commit template configurado. Usa 'git commit' (sem -m) para ver o template."

db/up: ## Inicia o PostgreSQL com Docker
	docker-compose up -d postgres

db/down: ## Para o PostgreSQL
	docker-compose down

db/logs: ## Mostra os logs do PostgreSQL
	docker-compose logs -f postgres

docker/build: ## Constrói a imagem Docker da app
	docker build -t gorm-crm:latest .
	@echo "✅ Imagem gorm-crm:latest construída"
	@docker images gorm-crm --format "Tamanho: {{.Size}}"

docker/up: ## Inicia toda a stack (app + postgres) com Docker
	docker-compose up -d
	@echo "✅ Stack a correr. API: http://localhost:8080"

docker/down: ## Para toda a stack
	docker-compose down

docker/logs: ## Mostra logs da app em tempo real
	docker-compose logs -f api

docker/ps: ## Mostra o estado dos containers
	docker-compose ps

version: ## Mostra a versão actual (git describe)
	@echo "Version:   $(VERSION)"
	@echo "Commit:    $(COMMIT)"
	@echo "BuildTime: $(BUILDTIME)"

release: ## Cria e faz push de uma tag de release (uso: make release TAG=v1.0.0)
	@[ "$(TAG)" ] || (echo "❌ Usa: make release TAG=v1.0.0"; exit 1)
	git tag -a $(TAG) -m "Release $(TAG)"
	git push origin $(TAG)
	@echo "✅ Tag $(TAG) criada e publicada — CD pipeline iniciado"

clean: ## Remove binários e ficheiros temporários
	@rm -rf bin/ coverage.out coverage.html
	@echo "🧹 Limpo"

.DEFAULT_GOAL := help
