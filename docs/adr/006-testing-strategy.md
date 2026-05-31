# ADR-006 — Estratégia de Testes

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 14 — Testes Automatizados

---

## Contexto

Em Go há várias abordagens para testes que envolvem infraestrutura (DB, Redis, etc.):

| Abordagem | Prós | Contras |
|-----------|------|---------|
| Mocks manuais | Rápidos, sem dependências | Podem divergir da implementação real |
| `sqlmock` | DB mockado | Falsa sensação de segurança — SQL real não testado |
| **Testcontainers** | DB real em Docker | Mais lento, requer Docker |
| Fixtures em memória | Rápidos | Não testa integração real |

## Decisão

**Abordagem híbrida por camada:**

| Camada | Abordagem | Biblioteca |
|--------|-----------|-----------|
| Service layer | Mocks manuais (implementam a interface) | `testify/assert` |
| Repository layer | Testcontainers — DB real | `testcontainers-go` |
| Handler layer | `httptest` com service mockado | `net/http/httptest` |
| E2E | Testcontainers + HTTP client real | `testcontainers-go` |

## Razões

1. **Unitários com mocks manuais** — o Repository Pattern (interfaces) já define o contrato; os mocks implementam essa interface e são verificados em compile-time
2. **Testcontainers para integração** — testa queries SQL reais, migrations reais, comportamento real do DB; evita surpresas em produção
3. **Sem `sqlmock`** — dá falsa segurança; uma query errada passa no mock mas falha no DB real
4. **Hierarquia clara** — mocks para velocidade, containers para confiança

## Estrutura

```
tests/
├── unit/
│   ├── contact_service_test.go    # usa MockContactRepository
│   └── auth_service_test.go       # usa MockUserRepository
├── integration/
│   ├── contact_repo_test.go       # usa testcontainers (PG real)
│   └── deal_repo_test.go
└── e2e/
    └── contact_flow_test.go       # HTTP client + testcontainers
```

## Convenções

- Table-driven tests para cobrir múltiplos casos num único teste
- `t.Parallel()` em todos os testes unitários
- Setup/teardown com `TestMain` para containers de integração
- Naming: `TestServiceName_MethodName_Scenario`

```go
func TestContactService_CreateContact_DuplicateEmail(t *testing.T) {
    // ...
}
```
