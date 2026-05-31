# 🎯 CHALLENGE — Módulo 14: Testes Automatizados

---

### Nível 1 — SpyPublisher em task.Service

Os testes de `task.Service` existentes verificam regras de negócio mas não verificam os eventos publicados. Adiciona um teste que confirma que o evento correcto é publicado.

Para isso, implementa um `SpyPublisher`:

```go
type SpyPublisher struct {
    Events []events.Event
}

func (s *SpyPublisher) Publish(event events.Event) {
    s.Events = append(s.Events, event)
}
```

Escreve um teste que:
1. Cria uma task
2. Verifica que o evento de criação foi publicado
3. Marca como `done`
4. Verifica que o evento correcto foi publicado

> **Nota:** Se `task.Service` não aceitar uma interface `Publisher`, esse é o exercício do M12 (SOLID/DIP) aplicado aqui.

---

### Nível 2 — Integration tests para lead.Repository

Replica o padrão do `tests/integration/contact_repository_test.go` para o domínio `lead`:

```go
// tests/integration/lead_repository_test.go
func TestLeadRepository_SaveAndFind(t *testing.T) { ... }
func TestLeadRepository_UpdateStatus(t *testing.T) { ... }
func TestLeadRepository_FindByContact(t *testing.T) { ... }
```

Atenção ao campo `Value float64` — verifica que vai e vem do PostgreSQL sem perda de precisão.

> **Pergunta:** Se `Lead.Value` fosse `valueobject.Money` (como no M13), o que mudarias no repositório? E nos testes?

---

### Nível 3 — E2E do pipeline de leads

Adiciona `tests/e2e/lead_api_test.go` que testa o fluxo completo de um lead:

```
POST /leads              → 201, lead criado com status "new"
PATCH /leads/:id/status  → 200, status actualizado para "contacted"
PATCH /leads/:id/status  → 422, transição inválida (new → qualified)
DELETE /leads/:id        → 204
GET /leads/:id           → 404
```

Usa o mesmo padrão `newTestApp` do `contact_api_test.go` — middleware de teste em vez de JWT.

---

## Perguntas de reflexão

1. **Unit vs Integration:** Tens um bug no índice único do email em PostgreSQL. Qual das três camadas de teste o detecta? Por quê?

2. **t.Parallel() e estado partilhado:** Se dois testes paralelos usarem o mesmo repositório em memória, o que acontece? Como resolves?

3. **Cobertura não é tudo:** `go test -cover` diz 80%. O que é que esse número não te diz? Qual é o teste mais valioso que podes escrever?

4. **TestMain:** Quando usarias `TestMain` em vez de `newTestDB` por teste? Qual é o trade-off entre um container por suite vs um por teste?

---

> Módulo seguinte: [branch-15-patterns](https://github.com/titi-byte-dev/gorm-crm/tree/branch-15-patterns) — Design Patterns: Strategy, Observer, Repository, Factory aplicados ao GoRM
