# 🚀 GoRM — Um CRM construído em Go

> **Curso de backend com didática de autoconstrução.**
> Cada branch Git = 1 módulo de aprendizagem. Do zero ao deploy, do Júnior ao Sénior.

---

## 📦 Módulo 01 — Setup & Estrutura Go

> **Branch:** `branch-01-setup` | **Nível:** 🟢 Júnior | **Duração:** ~3 dias

### O que vais aprender

- Como estruturar um projeto Go com o Standard Layout
- O que são módulos Go (`go.mod`, `go.sum`)
- Como criar um servidor HTTP com Fiber
- Como organizar middlewares e error handling global
- Makefile para produtividade no dia-a-dia

### O que foi construído neste módulo

- `GET /health` — endpoint de saúde da app
- Error handler global que mapeia erros de domínio para HTTP
- Middleware de logging estruturado
- Package `pkg/logger` com `slog` (stdlib Go 1.21+)
- Makefile com comandos: `run`, `build`, `test`, `lint`, `tidy`

### Contexto no GoRM

```mermaid
flowchart LR
    A["go mod init"] --> B["Standard Layout\ncmd/ internal/ pkg/"]
    B --> C["Fiber setup\ncmd/api/main.go"]
    C --> D["Middlewares\nlogger + recover"]
    D --> E["GET /health\n→ 200 OK"]
    E --> F["make run ✅\nlocalhost:8080"]
```

---

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Modules](https://img.shields.io/badge/Módulos-18-blue)](docs/)

---

## 📖 O que é este repositório?

**GoRM** é simultaneamente um **curso de backend em Go** e uma **aplicação CRM real e funcional**.

A ideia é simples: **aprendes construindo**. Cada branch representa uma etapa de aprendizagem — fazes `git checkout` e estás imediatamente no contexto certo, com código funcional, documentação e desafios práticos.

No final, tens um CRM completo deployado, testado e documentado.

---

## 🗺️ Mapa do Curso

```mermaid
flowchart TD
    START([🚀 START]) --> M01

    subgraph JUNIOR["🟢 NÍVEL JÚNIOR"]
        M01[📦 M01 · Setup and Estrutura Go]
        M02[🔤 M02 · Fundamentos Go]
        M03[🗄️ M03 · SQL e PostgreSQL]
        M04[🌿 M04 · Git Workflow]
        M05[🌐 M05 · REST API]
        M06[🔐 M06 · Autenticacao e Auth]
        M07[🏗️ M07 · Arquitetura MVC]
        M08[🐳 M08 · Docker]
        M01 --> M02 --> M03 --> M04 --> M05 --> M06 --> M07 --> M08
    end

    M08 --> JUNIOR_BADGE([🏆 Programador Junior])

    subgraph PLENO["🔵 NÍVEL PLENO"]
        M09[🍃 M09 · NoSQL e MongoDB]
        M10[✨ M10 · Clean Code]
        M11[🧩 M11 · OOP Avancado]
        M12[🏛️ M12 · SOLID]
        M13[🤸 M13 · Object Calisthenics]
        M14[🧪 M14 · Testes Automatizados]
        M15[🎨 M15 · Design Patterns]
        M09 --> M10 --> M11 --> M12 --> M13 --> M14 --> M15
    end

    JUNIOR_BADGE --> M09

    subgraph SENIOR["🟣 NÍVEL SÉNIOR"]
        M16[🔧 M16 · Refactoring]
        M17[⚡ M17 · Performance e Cache]
        M18[☁️ M18 · Cloud e CI/CD]
        M16 --> M17 --> M18
    end

    M15 --> M16
    M18 --> SENIOR_BADGE([🎓 Programador Senior])
```

---

## 🌿 Navegação por Branches

```bash
git clone https://github.com/titi-byte-dev/gorm-crm.git
cd gorm-crm

# Navega para qualquer módulo
git checkout branch-01-setup
git checkout branch-05-rest-api
git checkout branch-12-solid
```

| Branch | Módulo | Nível | Feature no GoRM |
|--------|--------|-------|-----------------|
| `branch-01-setup` | Setup e Estrutura | 🟢 Junior | Projeto inicializado, `GET /health` |
| `branch-02-go-fundamentos` | Fundamentos Go | 🟢 Junior | Domain models definidos |
| `branch-03-sql` | SQL e PostgreSQL | 🟢 Junior | CRUD de Contactos |
| `branch-04-git-workflow` | Git Workflow | 🟢 Junior | Branching strategy |
| `branch-05-rest-api` | REST API | 🟢 Junior | API REST completa |
| `branch-06-auth` | Autenticacao | 🟢 Junior | JWT + RBAC |
| `branch-07-mvc-layers` | Arquitetura MVC | 🟢 Junior | Camadas separadas |
| `branch-08-docker` | Docker | 🟢 **→ Junior** | App containerizada |
| `branch-09-nosql` | NoSQL e MongoDB | 🔵 Pleno | Activity logs |
| `branch-10-clean-code` | Clean Code | 🔵 Pleno | Codebase refatorada |
| `branch-11-oop` | OOP Avancado | 🔵 Pleno | Interfaces avancadas |
| `branch-12-solid` | SOLID | 🔵 Pleno | SOLID aplicado |
| `branch-13-calisthenics` | Object Calisthenics | 🔵 Pleno | Regras aplicadas |
| `branch-14-testes` | Testes | 🔵 Pleno | Unit + Integration + E2E |
| `branch-15-patterns` | Design Patterns | 🔵 Pleno | 10+ patterns aplicados |
| `branch-16-refactoring` | Refactoring | 🟣 Senior | Tecnicas avancadas |
| `branch-17-performance` | Performance e Cache | 🟣 Senior | Redis + Jobs async |
| `branch-18-cloud-cicd` | Cloud e CI/CD | 🟣 **→ Senior** | Deploy + Pipeline |

---

## 🏗️ Modelo de Dados

```mermaid
erDiagram
    USER {
        uuid id
        string name
        string email
        string password_hash
        enum role
        timestamp created_at
    }
    CONTACT {
        uuid id
        string name
        string email
        string phone
        string company
        uuid owner_id
    }
    LEAD {
        uuid id
        string title
        decimal value
        enum status
        uuid contact_id
        uuid owner_id
    }
    DEAL {
        uuid id
        string title
        decimal value
        enum stage
        uuid lead_id
        uuid contact_id
    }
    TASK {
        uuid id
        string title
        enum priority
        enum status
        uuid contact_id
        uuid deal_id
    }
    USER ||--o{ CONTACT : owns
    USER ||--o{ LEAD : owns
    CONTACT ||--o{ LEAD : generates
    LEAD ||--o| DEAL : converts
    CONTACT ||--o{ TASK : has
    DEAL ||--o{ TASK : has
```

---

## 🔄 Pipeline de Vendas

```mermaid
stateDiagram-v2
    [*] --> Contacto : Novo contacto
    Contacto --> Lead : Qualificado
    Lead --> Proposta : Interesse confirmado
    Proposta --> Negociacao : Proposta enviada
    Negociacao --> Ganho : Aceite
    Negociacao --> Perdido : Recusado
    Ganho --> [*]
    Perdido --> [*]
```

---

## 🏛️ Arquitetura Final

```mermaid
flowchart TD
    Client["Cliente HTTPS"] --> API

    subgraph APP["GoRM Application"]
        API["API Layer - Fiber"]
        Service["Service Layer - Business Logic"]
        Repo["Repository Layer - Data Access"]
        Worker["Workers - Goroutines Async"]
        Cache["Cache - Redis"]
        API --> Service
        Service --> Repo
        Service --> Worker
        Service <--> Cache
    end

    Repo --> PG[("PostgreSQL")]
    Repo --> MG[("MongoDB")]
    Worker --> EMAIL["SMTP - Notificacoes"]
```

---

## 🔐 Fluxo de Autenticacao

```mermaid
sequenceDiagram
    actor User
    participant API
    participant AuthService
    participant DB
    participant JWT

    User->>API: POST /auth/login
    API->>AuthService: Login(credentials)
    AuthService->>DB: SELECT user
    DB-->>AuthService: user record

    alt Password valida
        AuthService->>JWT: GenerateToken(userID, role)
        JWT-->>AuthService: signed token
        API-->>User: 200 OK com token
    else Password invalida
        API-->>User: 401 Unauthorized
    end

    User->>API: GET /contacts com Bearer token
    API->>JWT: ValidateToken()
    JWT-->>API: claims
    API-->>User: 200 OK com contacts
```

---

## 🧪 Piramide de Testes

```mermaid
flowchart TB
    E2E["E2E - 5 percent - Fluxos completos via HTTP"]
    INT["Integracao - 15 percent - Service e Repository com DB real"]
    UNIT["Unitarios - 80 percent - Service logic, Validacoes, Mappers"]
    UNIT --> INT --> E2E
    style UNIT fill:#22c55e,color:#fff
    style INT fill:#f59e0b,color:#fff
    style E2E fill:#ef4444,color:#fff
```

---

## ⚙️ CI/CD Pipeline

```mermaid
flowchart LR
    DEV["git push"] --> PR["Pull Request"]

    subgraph CI["GitHub Actions CI"]
        L["golangci-lint"] --> T["go test"] --> B["go build"] --> D["docker build"]
    end

    PR --> CI
    CI -->|pass| MERGE["Merge to main"]

    subgraph CD["Deploy Pipeline"]
        R["Push Registry"] --> S["Deploy Staging"] --> ST["Smoke Tests"] --> P["Deploy Prod"]
    end

    MERGE --> CD
```

---

## 📁 Estrutura de Pastas

```
gorm-crm/
├── cmd/api/main.go
├── internal/
│   ├── contact/        # Handler · Service · Repository · Model · DTO
│   ├── lead/
│   ├── deal/
│   ├── task/
│   ├── auth/           # JWT · Middleware · RBAC
│   └── shared/         # Errors · Events · Utils
├── pkg/
│   ├── database/
│   ├── cache/
│   └── logger/
├── migrations/
├── docs/               # Diagramas e ADRs
├── tests/
│   ├── unit/
│   ├── integration/
│   └── e2e/
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .github/workflows/ci.yml
```

---

## 🛠️ Stack

| Camada | Tecnologia |
|--------|-----------|
| Linguagem | Go 1.22+ |
| HTTP Framework | Fiber |
| ORM | GORM + golang-migrate |
| Base de dados | PostgreSQL |
| Logs / Historico | MongoDB |
| Cache | Redis |
| Auth | JWT (golang-jwt) |
| Containers | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Cloud | AWS / GCP |
| Testes | testify + testcontainers |

---

## 🚦 Como Comecar

```bash
git clone https://github.com/titi-byte-dev/gorm-crm.git
cd gorm-crm
git checkout branch-01-setup
cat README.md
make run
curl http://localhost:8080/health
```

---

## 📋 Checklist por Modulo

- [ ] README.md com objetivo claro
- [ ] Diagrama Mermaid de contexto
- [ ] Codigo Go funcional e testavel
- [ ] Testes (a partir do M03)
- [ ] CHALLENGE.md com exercicio pratico
- [ ] ADR se houve decisao de design
- [ ] git tag no final do modulo

---

## 📜 Licenca

MIT © [titi-byte-dev](https://github.com/titi-byte-dev)

---

> *"O melhor codigo e o codigo que tu proprio construiste, entendes e consegues explicar."*
