# 🎯 CHALLENGE — Módulo 07: Arquitectura MVC em Camadas

---

### Nível 1 — Testa com mock

**`MockContactRepository` e teste de email duplicado**

Cria um mock para `contact.Repository` em `tests/unit/contact_service_test.go` e testa:

```go
func TestContactService_Create_DuplicateEmail(t *testing.T) {
    svc := contact.NewService(newMockContactRepo(), bus)

    dto := contact.CreateContactDTO{Name: "Ana", Email: "ana@x.com", ...}
    svc.Create(ownerID, dto)           // primeiro — ok
    _, err := svc.Create(ownerID, dto) // segundo — deve falhar com ErrConflict

    if err == nil { t.Fatal("expected conflict error") }
}
```

Sem DB, sem servidor. Corre em < 1ms.

---

### Nível 2 — Experiência didáctica

**Quebra a arquitectura propositadamente e observa o impacto**

Muda o `task.Service.UpdateStatus` para aceitar `*fiber.Ctx` em vez de `task.Status`:

```go
// ❌ Propositadamente errado
func (s *Service) UpdateStatus(c *fiber.Ctx, id uuid.UUID) (*Task, error) {
    newStatus := Status(c.Query("status"))
    // ...
}
```

Tenta correr os testes unitários. O que acontece? Porquê?

Reverte a mudança e escreve no teu diário o que aprendeste.

---

### Nível 3 — Feature nova com arquitectura correcta

**`GET /api/v1/contacts/:id/tasks`**

Devolve as tasks associadas a um contacto.

A lógica certa:
1. `task.Repository` já tem `FindByContact(contactID)` — usar
2. Quem chama? O `contact.Handler`? O `task.Handler`? Porquê?
3. Precisa de um novo Service method? Ou acede ao repo directo?

Pensa antes de codificar. A resposta "certa" depende das tuas razões.

---

## Perguntas de reflexão

1. Se o Service precisar de enviar um email, onde deve estar esse código?
2. O que significa "testabilidade" e como as camadas a melhoram?
3. Há situações em que a separação em 3 camadas é excessiva? Quando simplificarias?

---

> Módulo seguinte: [branch-08-docker](https://github.com/titi-byte-dev/gorm-crm/tree/branch-08-docker) — Dockerfile, docker-compose e o ambiente de desenvolvimento completo
