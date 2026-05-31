# ADR-008 — Arquitectura em Camadas (Handler → Service → Repository)

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 07 — Arquitectura MVC

---

## Contexto

À medida que a codebase cresceu (M03→M06), o padrão Handler/Service/Repository
emergiu naturalmente. Este ADR documenta as regras e os **porquês** de cada camada.

## As Três Camadas

### Handler — "O que o mundo exterior vê"
**Responsabilidade única:** HTTP  
**Pode:** ler request, validar input, chamar service, serializar response  
**Não pode:** lógica de negócio, acesso direto ao DB, conhecer SQL  

```go
// ✅ Handler faz só HTTP
func (h *Handler) Create(c *fiber.Ctx) error {
    ownerID, _ := ctxutil.OwnerID(c)
    var dto CreateContactDTO
    c.BodyParser(&dto)
    validate.Check(dto)
    contact, err := h.svc.Create(ownerID, dto)  // delega ao service
    return response.Created(c, contact)
}

// ❌ Handler com lógica de negócio — errado
func (h *Handler) Create(c *fiber.Ctx) error {
    // verificar duplicados aqui? NÃO — isso é regra de negócio
    existing, _ := h.db.Where("email = ?", dto.Email).First(&rec)
}
```

### Service — "As regras do domínio"
**Responsabilidade única:** Lógica de negócio  
**Pode:** orquestrar repositórios, validar regras, emitir eventos  
**Não pode:** conhecer HTTP (fiber.Ctx), fazer queries SQL diretamente  

```go
// ✅ Service puro — sem HTTP, sem SQL
func (s *Service) Create(ownerID uuid.UUID, dto CreateContactDTO) (*Contact, error) {
    existing, _ := s.repo.FindByEmail(dto.Email)  // usa interface, não DB
    if existing != nil {
        return nil, fmt.Errorf("email exists: %w", ErrConflict)
    }
    contact := &Contact{...}
    saved, err := s.repo.Save(contact)
    s.bus.Publish(events.Event{Type: events.ContactCreated, Payload: saved})
    return saved, nil
}
```

### Repository — "Como os dados são guardados"
**Responsabilidade única:** Persistência  
**Pode:** SQL, GORM, queries, mapeamento DB ↔ domain  
**Não pode:** lógica de negócio, conhecer HTTP  

```go
// ✅ Repository puro — só DB
func (r *postgresRepository) FindByEmail(email string) (*Contact, error) {
    var rec contactRecord
    err := r.db.Where("email = ?", email).First(&rec).Error
    // mapeia rec → Contact (domain model)
    return recordToContact(rec), nil
}
```

## Por que esta separação?

| Cenário | Sem camadas | Com camadas |
|---------|-------------|-------------|
| Testar regras de negócio | Requer HTTP server + DB | Mock o repository, unit test rápido |
| Trocar PostgreSQL por MySQL | Alterar código espalhado | Só a implementação do Repository muda |
| Reutilizar lógica em CLI | Impossível — está misturada com HTTP | Chama o Service diretamente |
| Novo developer entende o código | Precisa de ler tudo | Sabe onde procurar pela camada |

## Regra de dependência

```
Handler → Service → Repository → DB
```

As dependências só apontam para baixo. O Service **nunca** importa o Handler.
O Repository **nunca** importa o Service. Violações desta regra são bugs de arquitectura.

## Consequências

- Cada ficheiro tem ~100 linhas porque tem UMA responsabilidade
- Testes unitários do Service não precisam de DB (usam mock do Repository)
- Adicionar um novo recurso (ex: Invoice) segue sempre o mesmo padrão
