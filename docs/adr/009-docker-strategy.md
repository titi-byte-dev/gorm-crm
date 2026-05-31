# ADR-009 — Estratégia Docker

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 08 — Docker

---

## Contexto

A app precisa de correr de forma consistente em qualquer ambiente:
máquina do developer, CI, staging e produção.

## Decisões

### Multi-stage Dockerfile

| Stage | Imagem base | Tamanho | Propósito |
|-------|-------------|---------|-----------|
| builder | golang:1.22-alpine | ~800MB | Compilar o binário |
| runtime | alpine:3.19 | ~5MB | Correr em produção |

Resultado: imagem de produção ~15MB (binário + alpine + ca-certs).

### Ordem das camadas no Dockerfile

```dockerfile
COPY go.mod go.sum ./    ← primeiro (muda raramente)
RUN go mod download      ← cacheia aqui
COPY . .                 ← depois (muda frequentemente)
RUN go build ...
```

O Docker invalida o cache a partir da primeira linha que muda.
Se copiarmos o código antes das dependências, cada mudança de código
invalida o cache das dependências — descarregando tudo de novo.

### Utilizador não-root

```dockerfile
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
```

Princípio de menor privilégio: o processo só tem os permissões mínimos.

### depends_on com healthcheck

```yaml
depends_on:
  postgres:
    condition: service_healthy
```

Garante que o postgres está realmente pronto (não só a iniciar)
antes de a api tentar ligar. Elimina race conditions no startup.

## Consequências

- `docker-compose up` inicia o ambiente completo em < 30s numa máquina limpa
- `make docker/build` produz uma imagem de ~15MB pronta para produção
- MongoDB e Redis estão comentados no compose — prontos para M09 e M17
