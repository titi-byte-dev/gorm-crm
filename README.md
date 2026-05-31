<!-- NAVIGATION BAR -->
<div align="center">

**[⬅️ M05 — REST API](https://github.com/titi-byte-dev/gorm-crm/tree/branch-05-rest-api)** &nbsp;|&nbsp;
`branch-06-auth` &nbsp;|&nbsp;
**[M07 — Arquitetura MVC ➡️](https://github.com/titi-byte-dev/gorm-crm/tree/branch-07-mvc-layers)**

`██████░░░░░░░░░░░░░░` Módulo **06 / 18** — Nível 🟢 Júnior

</div>

---

# 🔐 Módulo 06 — Autenticação & Autorização

[![CI](https://github.com/titi-byte-dev/gorm-crm/actions/workflows/ci.yml/badge.svg)](https://github.com/titi-byte-dev/gorm-crm/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![JWT](https://img.shields.io/badge/JWT-HS256-000000?style=flat&logo=jsonwebtokens)](.)
[![Módulo](https://img.shields.io/badge/Módulo-06%20%2F%2018-brightgreen)](.)

> **O que foi construído:** A API deixou de ser pública. JWT com access + refresh tokens, bcrypt para passwords, RBAC com roles (admin/manager/seller) e todas as rotas protegidas por middleware.

---

## 🎯 Objetivos de Aprendizagem

Ao terminar este módulo consegues:

- [ ] Explicar a diferença entre autenticação e autorização
- [ ] Implementar bcrypt e perceber porquê é lento por design
- [ ] Criar e validar JWTs com claims personalizados
- [ ] Usar middleware para proteger grupos de rotas
- [ ] Implementar RBAC com hierarquia de roles

---

## ⚡ Começa já

```bash
git checkout branch-06-auth
cp .env.example .env
# Edita .env e adiciona: JWT_SECRET=uma-chave-secreta-longa

docker-compose up -d postgres
make run
```

```bash
# 1. Criar conta
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Ana Silva","email":"ana@empresa.com","password":"segredo123","role":"seller"}'

# 2. Login — guarda o access_token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ana@empresa.com","password":"segredo123"}'

# 3. Aceder à API com o token
curl http://localhost:8080/api/v1/contacts \
  -H "Authorization: Bearer <access_token>"

# Sem token → 401 Unauthorized
curl http://localhost:8080/api/v1/contacts
```

---

## 🗺️ Como este módulo foi construído — commit a commit

> [!TIP]
> Corre `git log --oneline branch-05-rest-api..branch-06-auth` para ver todos os commits deste módulo com as suas explicações.

```mermaid
flowchart LR
    C1["chore: add JWT\n+ bcrypt deps\nPorquê estes packages?"]
    C2["feat: bcrypt\npassword hashing\nSalt · custo · timing"]
    C3["feat: JWT\ngeneration + validation\nHeader·Payload·Signature"]
    C4["feat: JWT\nmiddleware + RBAC\nProtected · RequireRole"]
    C5["feat: auth\nservice\nUser enumeration · refresh"]
    C6["feat: handler\n+ user repository\njson:\"-\" · /me endpoint"]
    C7["refactor: replace\nhardcoded ownerID\nctxutil · DRY"]
    C8["feat: wire auth\nprotect all routes\nroute groups"]
    C9["feat: tasks\nmigration\nPartial index"]

    C1-->C2-->C3-->C4-->C5-->C6-->C7-->C8-->C9
```

---

## 🔍 Conceitos-Chave

### JWT — Anatomia do token

```
eyJhbGciOiJIUzI1NiJ9 . eyJ1aWQiOiIxMjMiLCJyb2xlIjoic2VsbGVyIn0 . SflKxwRJS...
      HEADER                         PAYLOAD                          SIGNATURE
   (algoritmo)              (userID, role, exp, iat)              (HMAC do resto)
```

> [!WARNING]
> O Payload é apenas **Base64**, não encriptado. Qualquer um com o token consegue ler o conteúdo. **Nunca colocar passwords, dados de pagamento ou informação sensível no JWT.**

---

### bcrypt — Lento por design

<details>
<summary><strong>Ver: porquê 250ms é uma feature, não um bug</strong></summary>

```go
const bcryptCost = 12  // ~250ms numa máquina moderna

// Dois hashes da mesma password são SEMPRE diferentes
hash1, _ := HashPassword("segredo123")  // $2a$12$abc...xyz
hash2, _ := HashPassword("segredo123")  // $2a$12$def...uvw  ← diferente!

// A comparação é sempre em constant-time
CheckPassword("segredo123", hash1)  // true
```

**Porquê 250ms é bom:**
- Para um utilizador: imperceptível (login demora 300ms no total)
- Para um atacante com GPU: 10.000 tentativas/segundo × 250ms = inviável

**Porquê não SHA256:**
- SHA256 faz ~1 bilião de hashes/segundo numa GPU moderna
- 10 caracteres alfanuméricos: ~7 minutos para quebrar
- Com bcrypt cost 12: ~200 anos

</details>

---

### RBAC — Route Groups com middleware

<details>
<summary><strong>Ver: como as rotas ficam organizadas</strong></summary>

```go
v1 := app.Group("/api/v1")

// Públicas — sem middleware
auth.RegisterRoutes(v1, authSvc)

// Protegidas — Protected() corre antes de qualquer handler
protected := v1.Use(auth.Protected())
contact.RegisterRoutes(protected, contactSvc)

// Só admins — RequireRole corre depois de Protected()
adminOnly := protected.Use(auth.RequireRole(user.RoleAdmin))
// user.RegisterRoutes(adminOnly, userSvc)  ← Módulo 07
```

**Hierarquia:**
```
admin (3) → pode tudo
manager (2) → pode o que manager e seller podem
seller (1) → acesso básico
```

</details>

---

## 📁 Ficheiros deste módulo

<details>
<summary><strong>Ver ficheiros criados/modificados</strong></summary>

```
Criados:
├── internal/auth/
│   ├── password.go    ← bcrypt hash + check
│   ├── jwt.go         ← GenerateTokenPair + ValidateToken
│   ├── middleware.go  ← Protected() + RequireRole() + RBAC hierarchy
│   ├── service.go     ← Register, Login (user enumeration safe), Refresh
│   └── handler.go     ← /register /login /refresh /me
├── internal/user/
│   └── repository_pg.go
├── internal/shared/
│   └── ctxutil/ctxutil.go  ← OwnerID(c) helper partilhado
└── migrations/005_create_tasks.up/down.sql

Modificados:
├── internal/contact/handler.go  ← ownerID via ctxutil
├── internal/lead/handler.go     ← ownerID via ctxutil
├── internal/deal/handler.go     ← ownerID via ctxutil
└── cmd/api/main.go              ← auth wired + protected route group
```

</details>

---

## 🎯 Desafio

Ver [CHALLENGE.md](CHALLENGE.md)

- **Nível 1** — Testa user enumeration: tenta login com email que não existe vs password errada — as respostas são iguais?
- **Nível 2** — Implementa `PATCH /auth/password` para alterar password (requer token válido)
- **Nível 3** — Adiciona um endpoint só para admins e testa com um token de seller

---

## ✅ Checklist antes de avançar

- [ ] Fluxo completo testado: register → login → usar API → refresh
- [ ] Tentaste aceder a `/contacts` sem token — viste 401?
- [ ] Consegues explicar a diferença entre o access e o refresh token
- [ ] Entendes porquê o bcrypt é lento propositadamente

---

<!-- NAVIGATION BAR BOTTOM -->
<div align="center">

**[⬅️ M05 — REST API](https://github.com/titi-byte-dev/gorm-crm/tree/branch-05-rest-api)** &nbsp;|&nbsp;
`06 / 18` &nbsp;|&nbsp;
**[M07 — Arquitetura MVC ➡️](https://github.com/titi-byte-dev/gorm-crm/tree/branch-07-mvc-layers)**

</div>
