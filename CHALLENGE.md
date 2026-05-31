# 🎯 CHALLENGE — Módulo 15: Design Patterns

---

### Nível 1 — DealScorer (Strategy)

Replica o padrão `lead.Scorer` para o domínio `deal`:

```go
// internal/deal/scorer.go
type Scorer interface {
    Score(deal *Deal) int
}

type BasicDealScorer struct{}
// Critérios sugeridos:
// - StageNegotiation: +20, StageWon: +50, StageLost: -30
// - Value >= 50000: +40, Value >= 10000: +20

type UrgencyScorer struct {
    DaysUntilClose int // penaliza deals sem ClosedAt próximo
}
```

Integra no `deal.Service` com o mesmo padrão variadic.

> **Pergunta:** `DealScorer` e `LeadScorer` têm assinaturas idênticas mas tipos diferentes. Em Go, poderias ter uma única interface `Scorer[T any]`? Qual é o trade-off?

---

### Nível 2 — CachingContactRepo (Decorator)

Cria um segundo decorator que encadeia com o `ContactRepoLogger`:

```go
// pkg/decorator/contact_cache.go
type CachingContactRepo struct {
    inner  contact.Repository
    cache  sync.Map  // key: uuid.UUID, value: *contact.Contact
}

func (r *CachingContactRepo) FindByID(id uuid.UUID) (*contact.Contact, error) {
    if cached, ok := r.cache.Load(id); ok {
        return cached.(*contact.Contact), nil
    }
    result, err := r.inner.FindByID(id)
    if err == nil {
        r.cache.Store(id, result)
    }
    return result, err
}
```

Encadeia os dois decorators:
```go
repo := decorator.NewContactRepoLogger(
    decorator.NewCachingContactRepo(postgresRepo),
    logger,
)
```

> **Pergunta:** O `Save` e o `Update` devem invalidar o cache? E o `Delete`? Implementa a invalidação.

---

### Nível 3 — MaxContactsPerOwnerRule (Chain)

Adiciona uma regra ao pipeline de validação que limita o número de contactos por owner:

```go
// internal/contact/validation.go
type MaxContactsRule struct {
    Limit int
}

func (r MaxContactsRule) Validate(repo Reader, dto CreateContactDTO) error {
    // Reader precisa de um método Count(ownerID uuid.UUID) (int64, error)
    // Terás de adicionar Count à interface Reader
}
```

> **Atenção:** Para implementar esta regra, terás de estender a interface `Reader`. Isso é uma mudança breaking para todos os mocks que a implementam — como resolves?

---

## Perguntas de reflexão

1. **Strategy vs if/else:** Quando é que um `switch` simples é melhor que o Strategy pattern? Onde traças a linha?

2. **Decorator e testes:** `ContactRepoLogger` é mais fácil ou mais difícil de testar do que o repositório original? Porquê?

3. **Chain e ordem:** A ordem das regras na `Chain` importa? Dá um exemplo em que trocar `UniqueEmailRule` e `EmailDomainRule` produz comportamento diferente.

4. **Padrões e Go idiomático:** O GoF (Gang of Four) definiu 23 padrões. Go resolve alguns deles de forma diferente — `sync.Once` substitui Singleton, funções de primeira classe substituem muitos Command/Strategy simples. Que padrões GoF não fazem sentido em Go?

---

> Módulo seguinte: [branch-16-refactoring](https://github.com/titi-byte-dev/gorm-crm/tree/branch-16-refactoring) — Refactoring: identificar e eliminar code smells com técnicas sistemáticas
