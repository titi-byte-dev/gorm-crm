# 🎯 CHALLENGE — Módulo 11: OOP Avançado em Go

---

### Nível 1 — Label() para EntityType

O tipo `activitylog.EntityType` já tem constantes (`EntityContact`, `EntityLead`, etc.) mas ainda não tem `Label()` em português.

Implementa:
```go
func (e EntityType) String() string { ... }
func (e EntityType) Label() string  { ... }  // "Contacto", "Lead", "Negócio", etc.
```

Verifica que compila:
```bash
go build ./internal/activitylog/...
```

---

### Nível 2 — Interface Labeler

Cria uma interface comum para todos os tipos que têm `Label()`:

```go
// Onde colocas esta interface? internal/shared/? Porquê?
type Labeler interface {
    Label() string
}
```

Depois verifica em compile-time que os 4 tipos a satisfazem:
```go
var _ Labeler = lead.StatusNew
var _ Labeler = deal.StageProposal
var _ Labeler = task.StatusTodo
var _ Labeler = task.PriorityHigh
```

> **Dica:** Em Go, onde defines a interface importa — define-a onde ela é *usada*, não onde os tipos são definidos. Porquê?

---

### Nível 3 — WithTimeout no Bus

Adiciona uma nova opção ao `events.Bus` que define um timeout para o processamento de cada evento:

```go
func WithTimeout(d time.Duration) Option { ... }
```

Usa-a no `dispatch`:
```go
func (b *Bus) dispatch(ctx context.Context, event Event) {
    if b.timeout > 0 {
        var cancel context.CancelFunc
        ctx, cancel = context.WithTimeout(ctx, b.timeout)
        defer cancel()
    }
    // ... handlers
}
```

Verifica que:
1. Os callers existentes (`events.New(WithBufferSize(...), WithLogger(...))`) não precisam de mudar
2. `go build ./...` passa sem erros

---

## Perguntas de reflexão

1. **Embedding vs herança**: Se `pagination.Base` tiver um campo `Page` e `contact.Filters` também tiver um campo `Page`, o que acontece? Experimenta.

2. **Interface implícita**: Qual a vantagem de Go não exigir `implements`? Pensa num cenário onde queres que um tipo de uma biblioteca externa satisfaça uma interface tua.

3. **Functional Options**: O padrão tem um custo — cada opção é uma alocação de closure. Em código de alta performance (milhos de chamadas/segundo), isso pode importar. Como resolvias?

---

> Módulo seguinte: [branch-12-solid](https://github.com/titi-byte-dev/gorm-crm/tree/branch-12-solid) — SOLID: os 5 princípios aplicados ao GoRM
