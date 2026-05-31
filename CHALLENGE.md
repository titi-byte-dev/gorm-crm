# 🎯 CHALLENGE — Módulo 09: NoSQL & MongoDB

---

### Nível 1 — Observar o Observer em acção

```bash
make docker/up

# Regista e faz login
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin","email":"admin@crm.com","password":"segredo123","role":"admin"}'

TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@crm.com","password":"segredo123"}' | jq -r .access_token)

# Cria um contacto
curl -s -X POST http://localhost:8080/api/v1/contacts \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Ana Silva","email":"ana@empresa.com","company":"TechCorp"}'

# Espera 1 segundo e vê os logs
sleep 1
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/activity/me | jq .
```

Deves ver um log com `"action": "contact.created"`.

---

### Nível 2 — Query por entidade

Usando o ID do contacto criado acima:

```bash
CONTACT_ID="<uuid do contacto>"
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/activity/contact/$CONTACT_ID"
```

---

### Nível 3 — Inspecionar o MongoDB directamente

```bash
# Entra no shell do MongoDB
docker exec -it gorm-crm-mongo mongosh

# No mongosh:
use gorm_crm_logs
db.activity_logs.find().pretty()
db.activity_logs.getIndexes()
```

Verifica que os índices foram criados e que o TTL index existe.

---

## Perguntas de reflexão

1. O que acontece aos logs de actividade se pararmos o MongoDB e criarmos 5 contactos?
2. Porquê o `handleEvent` não propaga o erro se o `repo.Save` falhar?
3. Qual a diferença entre um TTL index no MongoDB e um cron job de limpeza?

---

> Módulo seguinte: [branch-10-clean-code](https://github.com/titi-byte-dev/gorm-crm/tree/branch-10-clean-code) — Clean Code Principles aplicados ao GoRM
