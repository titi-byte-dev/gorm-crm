# ADR-005 — Estratégia de Error Handling

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Todos os módulos

---

## Contexto

Go trata erros como valores — é uma diferença fundamental para quem vem de linguagens com exceções.

## Decisão

**Seguir o idioma Go rigorosamente:**

1. `error` é um valor de retorno, não uma exceção
2. Definir tipos de erro de domínio em `internal/shared/errors/`
3. Usar `fmt.Errorf("context: %w", err)` para wrapping com contexto
4. `panic` apenas para erros de programação (invariantes impossíveis), nunca para erros de runtime
5. Erros de HTTP mapeados em camada de handler, nunca em service/repository

```go
// ✅ Correto — erro como valor com contexto
func (s *ContactService) GetContact(id uuid.UUID) (*Contact, error) {
    contact, err := s.repo.FindByID(id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, fmt.Errorf("contact %s: %w", id, ErrNotFound)
        }
        return nil, fmt.Errorf("contact service get: %w", err)
    }
    return contact, nil
}

// ✅ Handler mapeia para HTTP
func (h *ContactHandler) Get(c *fiber.Ctx) error {
    contact, err := h.svc.GetContact(id)
    if errors.Is(err, shared.ErrNotFound) {
        return c.Status(404).JSON(ErrorResponse{Message: "contact not found"})
    }
    if err != nil {
        return c.Status(500).JSON(ErrorResponse{Message: "internal error"})
    }
    return c.JSON(contact)
}
```

## Tipos de erro de domínio

```go
// internal/shared/errors/errors.go
var (
    ErrNotFound      = errors.New("not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrConflict      = errors.New("conflict")
    ErrValidation    = errors.New("validation error")
)
```

## Razões

1. **Idiomático** — é a forma que a comunidade Go espera
2. **Explícito** — não há surpresas ocultas (exceções não apanhadas)
3. **Testável** — `errors.Is()` e `errors.As()` tornam os testes de erro simples
4. **Rastreável** — wrapping com contexto facilita o debugging em produção
