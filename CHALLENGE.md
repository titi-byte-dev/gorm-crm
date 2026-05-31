# 🎯 CHALLENGE — Módulo 05: REST API Completa

---

### Nível 1 — Obrigatório

**`GET /api/v1/contacts/:id/leads`**

Lista todos os leads de um contacto específico. Usa `lead.Repository.FindByContact()`.

```bash
curl http://localhost:8080/api/v1/contacts/<contact-id>/leads
# → [{ lead1 }, { lead2 }, ...]
```

---

### Nível 2 — Exploração

**`GET /api/v1/pipeline`**

Devolve um resumo do pipeline de vendas:

```json
{
  "proposal":    { "count": 5, "total_value": 25000 },
  "negotiation": { "count": 3, "total_value": 18000 },
  "won":         { "count": 12, "total_value": 67000 },
  "lost":        { "count": 4, "total_value": 9000 }
}
```

Dica: usa `db.Model(&dealRecord{}).Select("stage, count(*), sum(value)").Group("stage").Find(&result)`

---

### Nível 3 — Rate Limiting

Adiciona um middleware de rate limiting: máximo 100 requests por minuto por IP.

Fiber tem `middleware/limiter` built-in:

```go
app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
}))
```

Testa com um loop:
```bash
for i in {1..110}; do curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/health; done
# Deve ver 200 nos primeiros 100 e 429 nos seguintes
```

---

## Perguntas de reflexão

1. Porque é que `PATCH /leads/:id/status` é melhor que `PUT /leads/:id` para mudar estado?
2. O que é idempotência numa API REST? Quais dos teus endpoints são idempotentes?
3. Qual a diferença entre `400 Bad Request` e `422 Unprocessable Entity`?

---

> Módulo seguinte: [branch-06-auth](https://github.com/titi-byte-dev/gorm-crm/tree/branch-06-auth) — JWT, bcrypt, refresh tokens e RBAC
