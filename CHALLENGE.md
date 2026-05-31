# 🎯 CHALLENGE — Módulo 02: Fundamentos Go

---

## O que já tens

- Domain models completos: `Contact`, `Lead`, `Deal`, `Task`, `User`
- Repository interfaces definidas para cada domínio
- Event Bus com goroutines e channels
- Table-driven tests para `Lead.Status` e `Task.IsOverdue`

---

## Desafios

### Nível 1 — Obrigatório

**1.1 — Método `Contact.DisplayName()`**

Adiciona ao model `Contact` um método que devolve uma representação legível:

```go
// Se tiver empresa: "João Silva (Acme Corp)"
// Se não tiver: "João Silva"
func (c Contact) DisplayName() string { ... }
```

Escreve um table-driven test para cobrir ambos os casos.

---

**1.2 — Validação de email no model**

Adiciona um método `Contact.ValidateEmail() error` que verifica se o email tem formato válido usando apenas a stdlib Go (package `regexp` ou `strings`).

```go
contact := Contact{Email: "nao-e-email"}
err := contact.ValidateEmail()
// err != nil
```

---

### Nível 2 — Exploração

**2.1 — Implementa um `MockContactRepository`**

Cria um mock que implementa `contact.Repository` usando um `map` em memória:

```go
type MockContactRepository struct {
    contacts map[uuid.UUID]*contact.Contact
    // Como garantir que o mock satisfaz a interface em compile-time?
    // Dica: var _ contact.Repository = (*MockContactRepository)(nil)
}
```

Escreve um teste que usa o mock para testar uma função que recebe `contact.Repository`.

---

**2.2 — Adiciona `Deal.DurationDays()`**

Método que calcula quantos dias um deal esteve aberto (entre `CreatedAt` e `ClosedAt`).
Se ainda estiver aberto, calcula até ao momento atual.

```go
deal.DurationDays() // int
```

---

### Nível 3 — Investigação

**3.1 — Goroutine leak**

Cria um programa simples que demonstra um goroutine leak (goroutine que nunca termina).
Depois corrige-o usando `context.WithCancel`.

Usa `runtime.NumGoroutine()` para verificar antes e depois.

**3.2 — Interface satisfaction em compile-time**

Investiga este padrão e explica porque é útil:

```go
var _ contact.Repository = (*PostgreSQLContactRepository)(nil)
```

---

## Perguntas de reflexão

1. Porque é que Go usa interfaces implícitas em vez de `implements` explícito?
2. Qual a diferença entre `make(chan Event)` e `make(chan Event, 500)`? O que acontece se o channel estiver cheio?
3. Porque é que `Task.IsOverdue()` tem receiver de valor e `Filters.SetDefaults()` tem receiver de ponteiro?

---

> Módulo seguinte: [branch-03-sql](https://github.com/titi-byte-dev/gorm-crm/tree/branch-03-sql) — PostgreSQL, GORM e CRUD de Contactos
