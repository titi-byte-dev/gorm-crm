# 🎯 CHALLENGE — Módulo 01: Setup & Estrutura Go

> Cada módulo tem um desafio prático. Não há resposta errada — o objetivo é explorar.

---

## O que já tens

- Projeto Go inicializado com Fiber
- `GET /health` funcional
- Estrutura de pastas (`cmd/`, `internal/`, `pkg/`)
- Makefile com comandos comuns

---

## Desafios

### Nível 1 — Obrigatório

**1.1 — Adiciona um endpoint `GET /api/v1/version`**

Devolve informação sobre a versão da app:

```json
{
  "version": "0.1.0",
  "go_version": "go1.22",
  "build_time": "2026-06-01T10:00:00Z"
}
```

Dicas:
- O `go_version` pode ser obtido com o package `runtime`
- O `build_time` pode ser injetado em compile-time via `-ldflags`

---

**1.2 — Adiciona um middleware de `RequestID`**

Cada request deve ter um ID único no header de resposta:

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

Dicas:
- Fiber tem um middleware `requestid` built-in
- Ou usa o package `github.com/google/uuid` para gerar o ID manualmente

---

### Nível 2 — Exploração

**2.1 — Graceful Shutdown**

A app deve terminar de forma limpa quando recebe `SIGTERM` ou `SIGINT` (Ctrl+C), esperando que os requests em curso terminem antes de fechar.

Dicas:
- `os/signal` para capturar sinais
- `app.ShutdownWithTimeout(5 * time.Second)`

---

**2.2 — Endpoint `GET /api/v1/ping` com latência simulada**

Adiciona um parâmetro de query `?delay=500` que simula latência em milissegundos.
Observa o que aparece nos logs com o middleware de logger.

---

### Nível 3 — Investigação

**3.1 — Lê sobre o Standard Go Layout**

- Por que se usa `internal/` em vez de colocar tudo na raiz?
- Qual a diferença entre `internal/` e `pkg/`?
- O que é o `cmd/` e porquê ter uma pasta por executável?

Referência: [github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

---

## Como submeter o teu trabalho

```bash
# Cria uma branch pessoal a partir desta
git checkout -b meu-m01-challenge

# Faz as tuas alterações e commits
git add -p
git commit -m "feat: add version endpoint and request-id middleware"

# Compara com a solução
git diff branch-01-setup..meu-m01-challenge
```

---

## Perguntas de reflexão

1. O que acontece se correres dois processos Go na mesma porta? Experimenta.
2. Qual a diferença entre `log.Fatal()` e `panic()`? Quando usar cada um?
3. Porque é que o `Makefile` usa `.PHONY`? O que acontece sem isso?

---

> Módulo seguinte: [branch-02-go-fundamentos](https://github.com/titi-byte-dev/gorm-crm/tree/branch-02-go-fundamentos) — Domain Models, Interfaces e Goroutines
