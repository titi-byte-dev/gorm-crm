# 🎯 CHALLENGE — Módulo 08: Docker

---

### Nível 1 — Observar o impacto do multi-stage

```bash
# Constrói a imagem
make docker/build

# Compara com a imagem base de Go
docker images | grep -E "gorm-crm|golang"
```

A diferença deve ser ~800MB (golang) vs ~15MB (gorm-crm).

**Questão:** o que aconteceria se usasses `FROM golang:1.22-alpine` no stage final em vez de `FROM alpine:3.19`?

---

### Nível 2 — Observar o healthcheck em acção

```bash
make docker/up

# Para o postgres enquanto a api está a correr
docker-compose stop postgres

# Verifica o healthcheck
curl http://localhost:8080/health
# → HTTP 503 { "status": "degraded" }

# Reinicia o postgres
docker-compose start postgres

# Aguarda ~10s e verifica de novo
curl http://localhost:8080/health
# → HTTP 200 { "status": "ok" }
```

---

### Nível 3 — docker-compose.test.yml

Cria um `docker-compose.test.yml` que:
1. Inicia um PostgreSQL temporário (sem volume — dados não persistem)
2. Corre `go test ./tests/integration/...` contra esse postgres
3. Para tudo no final

```yaml
# docker-compose.test.yml
services:
  postgres-test:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: gorm_crm_test
      # ...
```

Adiciona ao `Makefile`:
```makefile
test/integration: ## Testes de integração com Docker
    docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

---

## Perguntas de reflexão

1. O que é um "layer" no Docker e como o cache funciona?
2. Qual a diferença entre `docker-compose stop` e `docker-compose down`?
3. Se adicionares um ficheiro `.env` ao `.dockerignore`, como passas as variáveis de ambiente em produção?

---

> 🏆 **Parabéns — completaste o Nível Júnior!**
>
> Módulo seguinte: [branch-09-nosql](https://github.com/titi-byte-dev/gorm-crm/tree/branch-09-nosql) — MongoDB para activity logs — início do Nível Pleno
