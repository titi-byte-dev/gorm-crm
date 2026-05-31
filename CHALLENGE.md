# 🎯 CHALLENGE — Módulo 06: Autenticação & Autorização

---

### Nível 1 — Segurança observável

**Testa user enumeration protection**

Faz dois pedidos:
```bash
# Email que não existe
curl -X POST /api/v1/auth/login -d '{"email":"nao-existe@x.com","password":"qualquer"}'

# Email que existe, password errada
curl -X POST /api/v1/auth/login -d '{"email":"ana@empresa.com","password":"errada"}'
```

As respostas devem ser **idênticas** — mesmo status, mesmo body. Se forem diferentes, a API revela quais emails existem.

---

### Nível 2 — Exploração

**`PATCH /api/v1/auth/password`** — alterar password

Requer token válido. Input:
```json
{ "current_password": "...", "new_password": "..." }
```

Lógica:
1. Extrair userID do token
2. Ir ao DB buscar o user
3. Verificar `current_password` com `CheckPassword`
4. Fazer hash da `new_password` e guardar

---

### Nível 3 — RBAC em prática

Adiciona um endpoint `GET /api/v1/admin/users` que lista todos os utilizadores.
Protege-o com `RequireRole(user.RoleAdmin)`.

Testa:
```bash
# Com token de seller → 403 Forbidden
# Com token de admin  → 200 OK com lista de users
```

---

## Perguntas de reflexão

1. O JWT expira mas o utilizador continua com o token — quando é que isso é um problema real?
2. O que é um "token blacklist" e quando precisarias de um?
3. Qual a diferença entre `401 Unauthorized` e `403 Forbidden`?

---

> Módulo seguinte: [branch-07-mvc-layers](https://github.com/titi-byte-dev/gorm-crm/tree/branch-07-mvc-layers) — Separação em camadas, interfaces e injeção de dependências
