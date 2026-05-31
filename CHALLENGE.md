# 🎯 CHALLENGE — Módulo 10: Clean Code Principles

---

### Nível 1 — Caça aos números mágicos

Encontra mais 2 números mágicos no codebase que ainda não foram nomeados:

```bash
# Procura por números literais no código Go
grep -rn "[0-9]\+" internal/ pkg/ cmd/ --include="*.go" | grep -v "_test.go" | grep -v "//.*[0-9]"
```

Para cada um que encontrares:
1. Percebe o que representa
2. Dá-lhe um nome descritivo
3. Define como constante no ficheiro adequado

---

### Nível 2 — Renomear em vez de comentar

Encontra um comentário no codebase que explica O QUÊ o código faz (não o PORQUÊ).
Tenta **renomear** uma variável, função ou tipo para tornar o comentário desnecessário.

Exemplo do tipo de coisa a procurar:
```go
// verifica se o token expirou
if time.Now().After(exp) { ... }

// Solução: extrair para função com nome expressivo
if tokenIsExpired(exp) { ... }
```

---

### Nível 3 — Early return num handler

Olha para os handlers em `internal/*/handler.go`. Encontra um que tenha um `if/else` onde o else podia ser evitado com early return ou extracção de função.

Refactora-o mantendo o comportamento exactamente igual.
Verifica com:
```bash
go test ./...
# comportamento não mudou
```

---

## Perguntas de reflexão

1. Há situações em que um comentário que explica O QUÊ é aceitável? (pensa em código de performance crítica com operações de bit)
2. Porque é que `type EntityType string` em vez de `type EntityType int`? Qual a diferença na prática?
3. O princípio "funções pequenas" tem um limite? Quando é que extrair mais funções piora o código?

---

> Módulo seguinte: [branch-11-oop](https://github.com/titi-byte-dev/gorm-crm/tree/branch-11-oop) — OOP Avançado: interfaces, composição e DRY/KISS/YAGNI
