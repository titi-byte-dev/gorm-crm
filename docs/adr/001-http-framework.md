# ADR-001 — Escolha do Framework HTTP

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 01 — Setup & Estrutura

---

## Contexto

Go tem várias opções para HTTP:

| Opção | Prós | Contras |
|-------|------|---------|
| `net/http` nativo | Zero dependências, idiomático | Verbose para rotas, sem middleware built-in |
| **Fiber** | Alta performance, API familiar (Express-like), boa DX | Não usa `net/http` handlers |
| Gin | Popular, maduro, usa `net/http` | API mais verbose que Fiber |
| Echo | Equilibrado, bom middleware | Menor comunidade |
| Chi | Leve, compatível com `net/http` | Menos features out-of-the-box |

## Decisão

**Usar Fiber** como framework HTTP principal.

## Razões

1. **Sintaxe familiar** — quem vem de Node.js/Express (muito comum) adapta-se mais rápido
2. **Alta performance** — construído sobre `fasthttp`, benchmarks consistentemente acima dos concorrentes
3. **DX excelente** — menos boilerplate para as tarefas mais comuns
4. **Didaticamente superior** — o estudante foca-se na lógica Go, não no boilerplate HTTP

## Consequências

- Os handlers Fiber não são compatíveis com `net/http` diretamente — no módulo avançado mostramos como adaptar se necessário
- O estudante aprende padrões (middleware, grupos de rotas, error handling) que são transferíveis para outros frameworks

## Alternativa considerada

`net/http` nativo seria introduzido como módulo bónus para demonstrar o que o framework abstrai.
