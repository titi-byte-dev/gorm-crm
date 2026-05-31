# ─────────────────────────────────────────────
# STAGE 1 — builder
# Imagem completa de Go apenas para compilar.
# Esta imagem (~800MB) nunca chega a produção.
# ─────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Instala dependências de sistema necessárias para CGO
# (o driver PostgreSQL pgx usa CGO para melhor performance)
RUN apk add --no-cache git

WORKDIR /app

# Copia go.mod e go.sum ANTES do código fonte.
# Porquê esta ordem?
# O Docker cacheia camadas. Se o código mudar mas as dependências não,
# esta camada fica em cache e `go mod download` não corre de novo.
# Poupar ~30s em cada build durante desenvolvimento.
COPY go.mod go.sum ./
RUN go mod download

# Agora copia o código fonte
COPY . .

# Compila o binário com otimizações para produção:
#   CGO_ENABLED=0  — binário estático, sem dependências de runtime C
#   GOOS=linux     — target OS explícito (mesmo que buildemos em Mac/Windows)
#   -ldflags "-s -w" — remove debug symbols e DWARF → imagem ~30% menor
#   -a             — força recompilação de todos os packages
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/gorm-crm \
    ./cmd/api/main.go

# ─────────────────────────────────────────────
# STAGE 2 — runtime
# Imagem mínima: só o binário compilado.
# scratch seria ainda menor mas não tem CA certificates.
# alpine tem ~5MB e inclui ca-certificates para HTTPS.
# ─────────────────────────────────────────────
FROM alpine:3.19 AS runtime

# Certificados SSL necessários para chamadas HTTPS externas (ex: SMTP, APIs)
RUN apk add --no-cache ca-certificates tzdata

# Cria utilizador não-root.
# Correr como root dentro do container é má prática de segurança:
# se o processo for comprometido, o atacante tem root no container.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

WORKDIR /app

# Copia APENAS o binário compilado do stage anterior
COPY --from=builder /app/bin/gorm-crm .

# Copia as migrations (necessárias em runtime para o comando migrate)
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# HEALTHCHECK — o Docker verifica se o container está saudável.
# Usado pelo orchestrator (Docker Swarm, Kubernetes) para saber
# quando reiniciar ou redirecionar tráfego.
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/health || exit 1

ENTRYPOINT ["./gorm-crm"]
