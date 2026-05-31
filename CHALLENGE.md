# 🎯 CHALLENGE — Módulo 16: Refactoring

---

### Nível 1 — Replace Conditional with Map em deal.MoveStage

`deal/service.go` ainda tem um if/else para selecionar o evento na `publishDealEvent`:

```go
evtType := events.DealLost
if d.Stage == StageWon {
    evtType = events.DealWon
}
```

Aplica o mesmo padrão de lookup table que usámos em `lead/service.go`:

```go
var dealEventByStage = map[Stage]events.EventType{
    StageWon:  events.DealWon,
    StageLost: events.DealLost,
}
```

> **Pergunta:** Se adicionarmos `StageNegotiation: events.DealNegotiating` à tabela, o código de publicação muda? O que é que isso te diz sobre o padrão?

---

### Nível 2 — Extract to Lookup Table em activitylog

`activitylog/service.go` tem `entityTypeFromEventType` com um switch de 8 casos:

```go
func entityTypeFromEventType(et events.EventType) EntityType {
    switch et {
    case events.ContactCreated, events.ContactUpdated, events.ContactDeleted:
        return EntityContact
    // ...
    }
}
```

Substitui por uma lookup table. Atenção: vários event types mapeiam para o mesmo EntityType.

```go
var entityByEventType = map[events.EventType]EntityType{
    events.ContactCreated: EntityContact,
    events.ContactUpdated: EntityContact,
    // ...
}
```

> **Pergunta:** A versão com switch agrupa visualmente os eventos por entidade. A versão com map perde essa estrutura. Como resolvias o trade-off entre legibilidade e extensibilidade?

---

### Nível 3 — Eliminar Duplicação em Filters.SetDefaults

`lead.Filters` e `deal.Filters` têm `SetDefaults()` quase idênticos:

```go
// lead/model.go         deal/model.go
func (f *Filters) SetDefaults() {
    if f.Page <= 0 { f.Page = 1 }
    if f.Limit <= 0 || f.Limit > 100 { f.Limit = 20 }
    if f.SortBy == "" { f.SortBy = "created_at" }
    if f.SortDir == "" { f.SortDir = "desc" }
}
```

Três abordagens possíveis — analisa cada uma:

1. **Embedding**: criar `shared.PaginationFilters` com `SetDefaults()` e embeder em `lead.Filters` e `deal.Filters`
2. **Função livre**: `shared.SetFilterDefaults(page, limit *int, sortBy, sortDir *string)`
3. **Deixar como está**: duplicação tolerável se os domínios divergirem

> **Pergunta:** Qual das três abordagens escolherias? Porquê? Em que condições mudarias de opinião?

---

## Perguntas de reflexão

1. **Refactoring vs Reescrita:** Qual é a diferença prática? Quando é que um "refactoring" se torna uma reescrita?

2. **Testes como rede de segurança:** Os unit tests de M14 permitiram fazer estes refactorings com confiança. O que aconteceria sem eles?

3. **Package-level vs local:** `leadTransitions` é declarada a nível de package. Que trade-offs tem essa decisão (visibilidade, testabilidade, thread-safety)?

4. **Mikado Method:** Se quisesses refactorizar `Filters.SetDefaults()` para partilhado, por onde começarias? Que dependências precisariam de mudar primeiro?

---

> Módulo seguinte: [branch-17-performance](https://github.com/titi-byte-dev/gorm-crm/tree/branch-17-performance) — Performance & Cache: índices, N+1, e cache com Redis
