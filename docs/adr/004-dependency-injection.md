# ADR-004 — Injeção de Dependências

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 07 (MVC) em diante

---

## Contexto

Em Go existem várias abordagens para DI:

| Abordagem | Biblioteca | Prós | Contras |
|-----------|-----------|------|---------|
| **Manual (constructors)** | Nenhuma | Zero magia, explícito, simples | Verbose em apps grandes |
| Wire | google/wire | Code generation, compile-time safe | Curva de aprendizagem |
| fx | uber-go/fx | Runtime DI, poderoso | Muito abstrato para iniciantes |

## Decisão

**Injeção manual (constructors) até ao módulo 15, Wire introduzido no módulo 16 como refactor.**

```go
// Manual DI — explícito e claro
func NewContactService(repo ContactRepository, bus EventBus) *ContactService {
    return &ContactService{repo: repo, bus: bus}
}

// Em main.go
db := database.NewPostgres(cfg)
repo := contact.NewPostgreSQLRepository(db)
bus := events.NewEventBus()
svc := contact.NewContactService(repo, bus)
handler := contact.NewHandler(svc)
```

## Razões

1. **O estudante percebe o problema antes de usar a ferramenta** — quando o `main.go` cresce, a necessidade de Wire torna-se óbvia
2. **Sem magia** — cada dependência é explícita e rastreável
3. **Testabilidade imediata** — passar um `MockRepository` no constructor é direto, sem framework
4. **Progressão natural** — a introdução de Wire no M16 (Refactoring) demonstra como a ferramenta resolve um problema real que o estudante já sentiu

## Consequências

- O `main.go` vai crescer durante os primeiros módulos — isso é intencional e serve de motivação para Wire
- A interface de cada serviço fica bem definida desde cedo (necessário para DI manual)
