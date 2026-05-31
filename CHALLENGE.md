# 🎯 CHALLENGE — Módulo 12: SOLID em Go

---

### Nível 1 — UniquePhoneRule (OCP)

Sem modificar `contact/service.go`, adiciona uma regra que verifica se o telefone já existe.

```go
// contact/rules.go — adiciona este tipo
type UniquePhoneRule struct{}

func (r UniquePhoneRule) Validate(repo Reader, dto CreateContactDTO) error {
    // como verificas se o telefone já existe?
    // Dica: podes precisar de um novo método no Reader...
    // Ou devolves sempre nil por enquanto e justificas porquê
}
```

Passa a nova regra no `main.go`:
```go
contact.NewService(repo, bus,
    contact.UniqueEmailRule{},
    contact.UniquePhoneRule{},
)
```

Verifica que `go build ./...` passa.

> **Reflexão:** Se `Reader` não tiver `FindByPhone`, tens duas opções:
> 1. Adicionar `FindByPhone` ao `Reader` — o que implica implementar no `PostgresRepository`
> 2. Deixar a regra sempre passar (YAGNI) e documentar o porquê
> Qual escolhes? Porquê?

---

### Nível 2 — SpySubscriber (LSP)

Cria um `SpySubscriber` em `pkg/testutil`:

```go
type SpySubscriber struct {
    Subscriptions map[events.EventType]int  // conta quantas vezes cada tipo foi subscrito
}

var _ events.Subscriber = (*SpySubscriber)(nil)  // LSP compile-time

func (s *SpySubscriber) Subscribe(et events.EventType, _ events.Handler) {
    // implementa
}
```

Usa-o para verificar que `activitylog.Service.RegisterHandlers` subscreve exactamente 9 tipos de eventos:
```go
spy := &testutil.SpySubscriber{Subscriptions: make(map[events.EventType]int)}
svc.RegisterHandlers(spy)
// verifica len(spy.Subscriptions) == 9
```

---

### Nível 3 — Rule no lead.Service (DIP + OCP)

Aplica o mesmo padrão `Rule` ao `lead.Service`.

Cria `internal/lead/rules.go` com uma `Rule` interface e um `ContactExistsRule` que verifica que o `ContactID` no `CreateLeadDTO` corresponde a um contacto existente.

```go
type Rule interface {
    Validate(dto CreateLeadDTO) error
}
```

> **Nota:** Esta regra precisa de acesso ao `contact.Reader`. Como passas essa dependência?
> Funcional options? Constructor injection? Campo na struct?

---

## Perguntas de reflexão

1. **SRP vs coesão:** Extrair `EventMapper` tornou o código mais fácil de entender ou mais fragmentado? Onde traças a linha entre "demasiado pequeno" e "responsabilidade única"?

2. **OCP e custo:** O padrão `Rule` adiciona indireção — há mais tipos, mais ficheiros. Quando vale a pena? Quando é over-engineering?

3. **LSP e interfaces largas:** Se `events.Publisher` tiver 10 métodos, é mais difícil escrever implementações de teste corretas. Como isso se relaciona com ISP?

---

> Módulo seguinte: [branch-13-calisthenics](https://github.com/titi-byte-dev/gorm-crm/tree/branch-13-calisthenics) — Object Calisthenics: 9 regras de disciplina de código
