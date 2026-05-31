# 🎯 CHALLENGE — Módulo 03: SQL & PostgreSQL

---

## Desafios

### Nível 1 — Obrigatório

**Adiciona `GET /api/v1/contacts/stats`**

Endpoint que devolve estatísticas dos contactos do owner:

```json
{
  "total": 42,
  "by_company": [
    { "company": "Acme",  "count": 12 },
    { "company": "Globo", "count": 8  }
  ]
}
```

Dica: usa `db.Model(&contactRecord{}).Select("company, count(*) as count").Group("company").Find(&result)`

---

### Nível 2 — Exploração

**Soft Delete**

Em vez de apagar o registo, adiciona um campo `deleted_at *time.Time`. O `DELETE /contacts/:id` preenche esse campo. O `GET /contacts` só devolve registos onde `deleted_at IS NULL`.

GORM suporta isto nativamente com `gorm.Model` — investiga como.

---

### Nível 3 — Transações

**Criar Lead automaticamente ao criar Contacto**

Quando se cria um Contacto, cria também um Lead associado na mesma transação — se o Lead falhar, o Contacto não é criado.

```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 1. criar contacto com tx
    // 2. criar lead com tx
    // se qualquer um falhar, ambos são revertidos
    return nil
})
```

---

## Perguntas de reflexão

1. Porque é que separamos `contactRecord` (GORM) de `Contact` (domain)?
2. O que acontece se correres a migration `002` sem ter corrido a `001`?
3. Qual a diferença entre `db.Save()` e `db.Updates()` no GORM?

---

> Módulo seguinte: [branch-04-git-workflow](https://github.com/titi-byte-dev/gorm-crm/tree/branch-04-git-workflow)
