# 🎯 CHALLENGE — Módulo 17: Performance & Cache

---

### Nível 1 — FindByIDs em lead.Repository

O `contact.Repository` ganhou `FindByIDs` para eliminar N+1. Faz o mesmo para `lead`:

```go
// internal/lead/model.go
type Repository interface {
    // ...
    FindByIDs(ids []uuid.UUID) ([]*Lead, error)
}

// internal/lead/repository_pg.go
func (r *postgresRepository) FindByIDs(ids []uuid.UUID) ([]*Lead, error) {
    if len(ids) == 0 {
        return nil, nil
    }
    var recs []leadRecord
    // GORM com WHERE id IN ?
}
```

Escreve um teste unitário com mock que verifica que `FindByIDs` com 3 IDs devolve exactamente os 3 leads correspondentes.

> **Pergunta:** O que acontece se `ids` tiver duplicados? O PostgreSQL devolve duplicados também. Deves deduplicar antes de enviar, ou aceitar duplicados e deixar o caller lidar?

---

### Nível 2 — Stats na Cache TTL

Estende `pkg/cache.TTL` com contadores atómicos de hits e misses:

```go
type TTL[K comparable, V any] struct {
    mu    sync.Map
    ttl   time.Duration
    hits  atomic.Int64
    misses atomic.Int64
}

func (c *TTL[K, V]) Stats() (hits, misses int64) {
    return c.hits.Load(), c.misses.Load()
}
```

Usa `sync/atomic` — não `mu sync.Mutex` — para não bloquear `Get` só para contar.

> **Pergunta:** `atomic.Int64` vs `sync.Mutex` para os contadores — qual é a diferença de performance? Em que caso usarias Mutex de qualquer forma?

---

### Nível 3 — CachingLeadRepo + composição completa

Cria `pkg/decorator/lead_cache.go` seguindo o mesmo padrão que `contact_cache.go`.

Depois, no `cmd/server` (ou num teste de integração), compõe a cadeia completa:

```go
leadRepo := decorator.NewLeadRepoLogger(
    decorator.NewCachingLeadRepo(lead.NewPostgresRepository(db), 5*time.Minute),
    logger,
)
```

> **Pergunta:** `CachingContactRepo` e `CachingLeadRepo` são quase idênticos. Poderias criar um único `CachingRepo[T any]` genérico? Qual é o obstáculo principal?

---

## Perguntas de reflexão

1. **Cache invalidation:** "There are only two hard problems in computer science: cache invalidation and naming things." O `Delete` invalida, o `Update` aquece. Mas e o `FindAll`? Se um contacto é actualizado, a lista pode estar desactualizada. Como resolves?

2. **Composite index ordem:** `(owner_id, stage)` é diferente de `(stage, owner_id)`. Numa query `WHERE owner_id = ? AND stage = ?` qual deles é mais eficiente? E se a query for só `WHERE stage = ?`?

3. **N+1 vs JOIN:** A solução de batch loading usa 2 queries separadas. Uma alternativa seria um JOIN. Quando escolherias JOIN em vez de batch loading em Go?

4. **TTL vs LRU:** A cache usa TTL (expiração por tempo). Uma cache LRU expira por uso. Em que cenário usarias LRU em vez de TTL neste CRM?

---

> Módulo seguinte: [branch-18-cicd](https://github.com/titi-byte-dev/gorm-crm/tree/branch-18-cicd) — Cloud & CI/CD: Dockerfile, GitHub Actions, deploy
