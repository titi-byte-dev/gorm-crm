# 🎯 CHALLENGE — Módulo 13: Object Calisthenics

---

### Nível 1 — Leads first-class collection (Regra 4)

Aplica o mesmo padrão `Contacts` ao domínio `lead`.

```go
// internal/lead/model.go
type Leads []*Lead

func (ls Leads) FilterByStatus(s Status) Leads { ... }
func (ls Leads) TotalValue() valueobject.Money  { ... }
func (ls Leads) IDs() []uuid.UUID               { ... }
```

Actualiza `Reader.FindAll` e `Service.List` para devolver `Leads`.

> **Pergunta:** `TotalValue()` itera sobre todos os leads para somar. Se tiveres 10.000 leads em memória, qual é o impacto? Quando é que isto é um problema?

---

### Nível 2 — Email value object (Regra 3)

Cria `pkg/valueobject/email.go`:

```go
type Email string

func ParseEmail(s string) (Email, error) {
    // validações mínimas: não vazio, contém @, lowercase
}

func (e Email) String() string  { return string(e) }
func (e Email) Domain() string  { /* parte depois do @ */ }
func (e Email) IsEmpty() bool   { return e == "" }
```

> **Nota GORM:** `type Email string` funciona directamente — GORM armazena como VARCHAR.
> Não precisas de implementar `Scan`/`Value` — o tipo base é suficiente.

Aplica a `Contact.Email` e actualiza os conversores no repositório.

---

### Nível 3 — Caça aos níveis de indentação

Procura no codebase métodos com 2+ níveis de indentação:

```bash
# Padrão: dois ou mais tabs/espaços seguidos de "for" ou "if"
grep -rn "		if\|		for" internal/ --include="*.go"
```

Para cada ocorrência que encontrares:
1. Identifica a regra violada (1? 2? ambas?)
2. Refactoriza aplicando `slices.Contains`, extracção para função, ou guard clause
3. Verifica que `go build ./...` continua a passar

---

## Perguntas de reflexão

1. **Regra 3 vs. over-engineering:** `Money` faz sentido. Mas `type ContactID uuid.UUID`? Onde traças a linha entre "envolve o primitivo" e "indireção desnecessária"?

2. **Regra 4 e performance:** `Contacts.FilterByCompany()` cria um novo slice. Em Go, slices são leves — mas se fizeres 3 filtros encadeados, crias 3 slices intermédios. Como resolvias com `iter.Seq` (Go 1.23+)?

3. **Regra 6 e Go idiomático:** A comunidade Go tem convenções fortes para nomes curtos (`err`, `ctx`, `c`). Onde é que a Regra 6 entra em conflito com os idiomas Go? Como resolves o conflito?

---

> Módulo seguinte: [branch-14-tests](https://github.com/titi-byte-dev/gorm-crm/tree/branch-14-tests) — Testes Automatizados: unitários, integração com testcontainers, e2e
