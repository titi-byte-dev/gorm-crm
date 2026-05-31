# 📘 PRD + SSD — GoRM: Um CRM construído em Go

> **Product Requirements Document + Study & System Design**
> Versão: 1.0 | Data: 31/05/2026

---

## Índice

1. [Visão do Produto](#1-visão-do-produto)
2. [Objetivos](#2-objetivos)
3. [Público-Alvo](#3-público-alvo)
4. [Princípios de Design do Curso](#4-princípios-de-design-do-curso)
5. [App Âncora — GoRM CRM](#5-app-âncora--gorm-crm)
6. [Estrutura de Branches](#6-estrutura-de-branches)
7. [Arquitetura do Sistema](#7-arquitetura-do-sistema)
8. [Fluxos Detalhados](#8-fluxos-detalhados)
9. [Design Patterns Aplicados](#9-design-patterns-aplicados)
10. [Pirâmide de Testes](#10-pirâmide-de-testes)
11. [Pipeline CI/CD](#11-pipeline-cicd)
12. [Plano de Módulos](#12-plano-de-módulos)
13. [Roadmap de Implementação](#13-roadmap-de-implementação)
14. [Decisões de Design (ADRs)](#14-decisões-de-design-adrs)
15. [Checklist por Módulo](#15-checklist-por-módulo)

---

## 1. Visão do Produto

**GoRM** é simultaneamente um curso de backend em Go e uma aplicação real funcional.

O estudante aprende construindo — cada módulo adiciona uma feature ao CRM, e no final do percurso tem um produto deployado, testado e documentado.

> O repositório Git **é** o curso. Cada branch representa uma etapa de aprendizagem isolada, navegável e cumulativa.

---

## 2. Objetivos

| # | Objetivo | Métrica de sucesso |
|---|----------|--------------------|
| 1 | **Primário** — O estudante constrói um CRM completo em Go do zero ao deploy | App funcional na branch-18 |
| 2 | **Secundário** — O repositório serve como material reutilizável para ensinar outros | README + docs auto-suficientes por branch |
| 3 | **Regra 80/20** — Cada módulo cobre o que aparece em 80% dos projetos reais | Sem teoria sem aplicação prática |

---

## 3. Público-Alvo

| Perfil | Descrição |
|--------|-----------|
| **Primário** | Programador experiente noutra linguagem, que quer aprender Go e consolidar arquitetura backend |
| **Secundário** | Futuros alunos (qualquer nível, graças aos módulos didáticos progressivos) |

**Pré-requisitos assumidos:**
- Experiência em pelo menos uma linguagem de programação
- Familiaridade com conceitos de programação orientada a objetos
- Conhecimento básico de linha de comandos

---

## 4. Princípios de Design do Curso

```mermaid
flowchart LR
    ROOT(["⚙️ GoRM Course"])

    ROOT --> A["🔨 Autoconstrução"]
    A --> A1["Aprende fazendo"]
    A --> A2["Cada branch é uma feature real"]
    A --> A3["Código acumula ao longo do curso"]

    ROOT --> B["📊 80/20"]
    B --> B1["Foco no que é mais pedido"]
    B --> B2["Sem teoria sem aplicação"]
    B --> B3["Casos reais de mercado"]

    ROOT --> C["🧭 Navegabilidade"]
    C --> C1["Git como estrutura do curso"]
    C --> C2["Cada branch isolada e funcional"]
    C --> C3["git checkout é o índice"]

    ROOT --> D["♻️ Reutilizabilidade"]
    D --> D1["Documenta enquanto aprende"]
    D --> D2["README por branch"]
    D --> D3["Diagramas sempre presentes"]

    ROOT --> E["📈 Progressividade"]
    E --> E1["Do simples ao complexo"]
    E --> E2["Júnior → Pleno → Sénior"]
    E --> E3["Complexidade incremental real"]

    style ROOT fill:#1e40af,color:#fff,stroke:#1e3a8a
    style A fill:#065f46,color:#fff,stroke:#064e3b
    style B fill:#92400e,color:#fff,stroke:#78350f
    style C fill:#1e3a8a,color:#fff,stroke:#1e40af
    style D fill:#4c1d95,color:#fff,stroke:#3b0764
    style E fill:#831843,color:#fff,stroke:#500724
```

---

## 5. App Âncora — GoRM CRM

### 5.1 Porquê um CRM?

Um CRM cobre naturalmente os conceitos que aparecem em quase todos os sistemas backend reais:

| Conceito técnico | Presente no GoRM |
|------------------|-----------------|
| CRUD completo | Contactos, Leads, Deals, Tasks |
| Relações entre entidades | Cliente → Leads → Negócios → Tarefas |
| Autenticação + Roles | Vendedor, Manager, Admin |
| Estados e transições | Pipeline de vendas (Kanban) |
| Jobs assíncronos | Emails automáticos, follow-ups |
| Queries complexas | Funil de vendas, relatórios |
| Search full-text | Pesquisa de contactos |
| Multi-tenancy (avançado) | Empresas diferentes |

### 5.2 Modelo de Dados

```mermaid
erDiagram
    USER {
        uuid id PK
        string name
        string email
        string password_hash
        enum role "admin | manager | seller"
        timestamp created_at
        timestamp updated_at
    }
    CONTACT {
        uuid id PK
        string name
        string email
        string phone
        string company
        string notes
        uuid owner_id FK
        timestamp created_at
        timestamp updated_at
    }
    LEAD {
        uuid id PK
        string title
        decimal value
        enum status "new | contacted | qualified | lost"
        uuid contact_id FK
        uuid owner_id FK
        timestamp created_at
        timestamp updated_at
    }
    DEAL {
        uuid id PK
        string title
        decimal value
        enum stage "proposal | negotiation | won | lost"
        uuid lead_id FK
        uuid contact_id FK
        uuid owner_id FK
        timestamp closed_at
        timestamp created_at
    }
    TASK {
        uuid id PK
        string title
        string description
        enum priority "low | medium | high | urgent"
        enum status "todo | in_progress | done | cancelled"
        uuid assigned_to FK
        uuid contact_id FK
        uuid deal_id FK
        timestamp due_date
        timestamp created_at
    }
    ACTIVITY_LOG {
        uuid id PK
        string action
        json payload
        uuid user_id FK
        uuid entity_id
        string entity_type
        timestamp created_at
    }

    USER ||--o{ CONTACT : "owns"
    USER ||--o{ LEAD : "owns"
    USER ||--o{ DEAL : "owns"
    USER ||--o{ TASK : "assigned to"
    CONTACT ||--o{ LEAD : "generates"
    LEAD ||--o| DEAL : "converts to"
    CONTACT ||--o{ TASK : "has"
    DEAL ||--o{ TASK : "has"
    USER ||--o{ ACTIVITY_LOG : "generates"
```

### 5.3 Pipeline de Vendas

```mermaid
stateDiagram-v2
    [*] --> Contacto : Novo contacto adicionado
    Contacto --> Lead : Qualificado pelo vendedor
    Lead --> Proposta : Interesse confirmado
    Proposta --> Negociacao : Proposta enviada
    Negociacao --> Ganho : Aceite ✅
    Negociacao --> Perdido : Recusado ❌
    Ganho --> [*]
    Perdido --> [*]
    Lead --> Perdido : Desqualificado
    Contacto --> Perdido : Sem interesse
```

### 5.4 Roles e Permissões

```mermaid
flowchart TD
    subgraph ROLES["Sistema de Roles"]
        ADMIN["👑 Admin\n─────────\n• Gestão de utilizadores\n• Todos os contactos\n• Relatórios globais\n• Configurações"]
        MANAGER["👔 Manager\n─────────\n• Ver todos os deals\n• Reatribuir leads\n• Relatórios da equipa\n• Aprovar negócios"]
        SELLER["🧑‍💼 Seller\n─────────\n• Gerir próprios contactos\n• Criar leads e deals\n• Gerir próprias tasks\n• Ver pipeline pessoal"]
    end

    ADMIN --> MANAGER --> SELLER
```

---

## 6. Estrutura de Branches

### 6.1 Visão Geral

```mermaid
gitGraph
   commit id: "init: repo setup + README"

   branch branch-01-setup
   checkout branch-01-setup
   commit id: "feat: standard go layout"
   commit id: "feat: GET /health endpoint"
   commit id: "docs: module README + challenge"
   checkout main
   merge branch-01-setup tag: "v0.1.0"

   branch branch-02-go-fundamentos
   checkout branch-02-go-fundamentos
   commit id: "feat: domain models"
   commit id: "feat: interfaces definidas"
   commit id: "feat: goroutines intro"
   checkout main
   merge branch-02-go-fundamentos tag: "v0.2.0"

   branch branch-03-sql
   checkout branch-03-sql
   commit id: "feat: postgresql + gorm"
   commit id: "feat: migrations setup"
   commit id: "feat: contacts CRUD"
   checkout main
   merge branch-03-sql tag: "v0.3.0"

   branch branch-05-rest-api
   checkout branch-05-rest-api
   commit id: "feat: fiber setup"
   commit id: "feat: full REST routes"
   commit id: "feat: middlewares"
   checkout main
   merge branch-05-rest-api tag: "v0.5.0"

   branch branch-08-docker
   checkout branch-08-docker
   commit id: "feat: Dockerfile"
   commit id: "feat: docker-compose"
   checkout main
   merge branch-08-docker tag: "v0.8.0 🏆 JUNIOR"

   branch branch-15-patterns
   checkout branch-15-patterns
   commit id: "feat: repository pattern"
   commit id: "feat: observer events"
   commit id: "feat: factory + builder"
   checkout main
   merge branch-15-patterns tag: "v1.5.0 🎯 PLENO"

   branch branch-18-cloud-cicd
   checkout branch-18-cloud-cicd
   commit id: "feat: github actions CI"
   commit id: "feat: deploy pipeline"
   commit id: "feat: production config"
   checkout main
   merge branch-18-cloud-cicd tag: "v2.0.0 🎓 SENIOR"
```

### 6.2 Tabela de Branches

| Branch | Módulo | Nível | Feature adicionada ao GoRM |
|--------|--------|-------|---------------------------|
| `branch-01-setup` | Setup & Estrutura | 🟢 Júnior | Projeto Go, estrutura de pastas, `GET /health` |
| `branch-02-go-fundamentos` | Fundamentos Go | 🟢 Júnior | Domain models, interfaces de repositório |
| `branch-03-sql` | SQL & PostgreSQL | 🟢 Júnior | CRUD de Contactos com PostgreSQL + GORM |
| `branch-04-git-workflow` | Git Workflow | 🟢 Júnior | Branching strategy, conventional commits |
| `branch-05-rest-api` | REST API | 🟢 Júnior | API REST completa (Fiber), middlewares |
| `branch-06-auth` | Autenticação & Auth | 🟢 Júnior | JWT, RBAC (admin/manager/seller) |
| `branch-07-mvc-layers` | Arquitetura MVC | 🟢 Júnior | Handler/Service/Repository separados |
| `branch-08-docker` | Docker | 🟢 **→ Júnior** | Dockerfile, docker-compose, ambiente local |
| `branch-09-nosql` | NoSQL & MongoDB | 🔵 Pleno | Activity logs em MongoDB |
| `branch-10-clean-code` | Clean Code | 🔵 Pleno | Refactor com princípios Clean Code |
| `branch-11-oop` | OOP Avançado | 🔵 Pleno | Interfaces avançadas, DRY/KISS/YAGNI |
| `branch-12-solid` | SOLID | 🔵 Pleno | SOLID aplicado ao GoRM |
| `branch-13-calisthenics` | Object Calisthenics | 🔵 Pleno | 9 regras aplicadas ao código |
| `branch-14-testes` | Testes Automatizados | 🔵 Pleno | Unit + Integration + E2E |
| `branch-15-patterns` | Design Patterns | 🔵 Pleno | Repository, Observer, Factory, Strategy |
| `branch-16-refactoring` | Técnicas de Refactoring | 🟣 Sénior | Refactor com casos reais do CRM |
| `branch-17-performance` | Performance & Cache | 🟣 Sénior | Redis, jobs assíncronos, CDN |
| `branch-18-cloud-cicd` | Cloud & CI/CD | 🟣 **→ Sénior** | Deploy AWS/GCP, GitHub Actions |

---

## 7. Arquitetura do Sistema

### 7.1 Visão de Alto Nível (C4 — Context)

```mermaid
flowchart TD
    USER["👤 Utilizador\nVendedor · Manager · Admin"]

    subgraph GORM_SYSTEM["GoRM CRM"]
        API["🌐 API Gateway\nGo + Fiber\nREST API · Auth · Routing"]
        SERVICE["⚙️ Business Layer\nGo\nRegras de negócio · Validações"]
        REPO["🗄️ Repository Layer\nGo\nAbstração de acesso a dados"]
        WORKER["⚡ Background Workers\nGoroutines\nJobs · Emails · Notificações"]
        CACHE["🔴 Cache Layer\nRedis\nSessões · Queries frequentes"]
    end

    PG[("🐘 PostgreSQL\nDados relacionais")]
    MONGO[("🍃 MongoDB\nActivity logs")]
    SMTP["📧 SMTP\nEmail Service"]
    CLOUD["☁️ Cloud Provider\nAWS / GCP"]

    USER -->|"HTTPS / JSON"| API
    API --> SERVICE
    SERVICE --> REPO
    SERVICE --> WORKER
    SERVICE <--> CACHE
    REPO --> PG
    REPO --> MONGO
    WORKER --> SMTP
    API --> CLOUD

    style GORM_SYSTEM fill:#f0f9ff,stroke:#0284c7
```

### 7.2 Estrutura Interna por Camadas

```mermaid
flowchart LR
    subgraph HTTP["🌐 HTTP Layer"]
        Router["Router\n(Fiber)"]
        Handler["Handlers\nContactHandler\nLeadHandler\nDealHandler"]
        Middleware["Middlewares\nAuth · Logger · CORS · RateLimit"]
    end

    subgraph SERVICE["⚙️ Service Layer"]
        ContactSvc["ContactService"]
        LeadSvc["LeadService"]
        DealSvc["DealService"]
        TaskSvc["TaskService"]
        AuthSvc["AuthService"]
    end

    subgraph REPO["🗄️ Repository Layer"]
        ContactRepo["ContactRepository\n(interface)"]
        PGContactRepo["PostgreSQLContactRepo\n(implementação)"]
        MockContactRepo["MockContactRepo\n(testes)"]
    end

    subgraph INFRA["🔧 Infrastructure"]
        PG[("PostgreSQL")]
        Redis[("Redis")]
        Mongo[("MongoDB")]
    end

    Router --> Middleware --> Handler
    Handler --> ContactSvc
    Handler --> LeadSvc
    ContactSvc --> ContactRepo
    ContactRepo --> PGContactRepo
    PGContactRepo --> PG
    ContactSvc --> Redis
    LeadSvc --> Mongo
```

### 7.3 Estrutura de Pastas (Standard Go Layout)

```
gorm-crm/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point — wire everything together
│
├── internal/                        # Código privado da aplicação
│   ├── contact/
│   │   ├── handler.go               # HTTP handlers (recebe req, devolve resp)
│   │   ├── service.go               # Business logic (regras de negócio)
│   │   ├── repository.go            # Interface do repositório
│   │   ├── repository_pg.go         # Implementação PostgreSQL
│   │   ├── model.go                 # Domain model (struct Contact)
│   │   └── dto.go                   # Request/Response DTOs
│   ├── lead/                        # mesma estrutura
│   ├── deal/                        # mesma estrutura
│   ├── task/                        # mesma estrutura
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── jwt.go                   # Token generation + validation
│   │   └── middleware.go            # Auth middleware
│   └── shared/
│       ├── middleware/              # CORS, logger, rate limit
│       ├── errors/                  # Error types + handlers
│       ├── events/                  # Event bus (channels)
│       └── pagination/              # Query params helpers
│
├── pkg/                             # Código reutilizável (pode ser importado externamente)
│   ├── database/
│   │   ├── postgres.go              # PostgreSQL connection
│   │   └── mongodb.go               # MongoDB connection
│   ├── cache/
│   │   └── redis.go                 # Redis client
│   └── logger/
│       └── logger.go                # Structured logging (zerolog/zap)
│
├── migrations/                      # SQL migrations (golang-migrate)
│   ├── 001_create_users.up.sql
│   ├── 001_create_users.down.sql
│   └── ...
│
├── docs/                            # Documentação e diagramas
│   ├── PRD-SSD.md                   # Este documento
│   ├── adr/                         # Architecture Decision Records
│   └── modules/                     # Docs detalhadas por módulo
│
├── tests/
│   ├── unit/                        # Testes unitários (mocks)
│   ├── integration/                 # Testes com DB real (testcontainers)
│   └── e2e/                         # Testes end-to-end via HTTP
│
├── docker-compose.yml               # PostgreSQL + MongoDB + Redis local
├── docker-compose.test.yml          # Stack para testes CI
├── Dockerfile                       # Multi-stage build
├── Makefile                         # make run | test | migrate | lint
├── .env.example                     # Variáveis de ambiente necessárias
├── .github/
│   └── workflows/
│       ├── ci.yml                   # CI pipeline
│       └── deploy.yml               # CD pipeline
└── README.md
```

---

## 8. Fluxos Detalhados

### 8.1 Fluxo de Autenticação (Módulo 06)

```mermaid
sequenceDiagram
    actor User
    participant API as API (Fiber)
    participant Auth as AuthService
    participant DB as PostgreSQL
    participant JWT as JWT

    User->>API: POST /auth/login {email, password}
    API->>Auth: Login(credentials)
    Auth->>DB: SELECT * FROM users WHERE email = ?
    DB-->>Auth: user record

    alt Password válida
        Auth->>Auth: bcrypt.CompareHash(hash, password)
        Auth->>JWT: GenerateToken(userID, role, exp)
        JWT-->>Auth: signed JWT token
        Auth->>JWT: GenerateRefreshToken(userID)
        JWT-->>Auth: refresh token
        Auth-->>API: TokenPair{token, refresh}
        API-->>User: 200 OK {token, refresh_token, expires_in}
    else Password inválida
        Auth-->>API: ErrInvalidCredentials
        API-->>User: 401 Unauthorized
    end

    Note over User,JWT: Requests autenticados

    User->>API: GET /contacts [Authorization: Bearer <token>]
    API->>JWT: ValidateToken(token)
    alt Token válido
        JWT-->>API: Claims{userID, role, exp}
        API->>DB: SELECT * FROM contacts WHERE owner_id = userID
        DB-->>API: []Contact
        API-->>User: 200 OK {contacts, pagination}
    else Token expirado
        API-->>User: 401 Unauthorized {error: "token expired"}
        Note over User: Usar refresh token para renovar
    end
```

### 8.2 Fluxo CRUD Completo — Repository Pattern (Módulo 07 + 15)

```mermaid
flowchart TD
    REQ["HTTP Request\nPOST /contacts\n{name, email, company}"] --> HANDLER

    subgraph HANDLER["Handler Layer"]
        H1["ContactHandler.Create()"]
        H2["Bind + Validate DTO"]
        H1 --> H2
    end

    subgraph SERVICE["Service Layer"]
        S1["ContactService.CreateContact(dto)"]
        S2["Verificar duplicados\n(email único)"]
        S3["Mapear DTO → Entity"]
        S4["Emitir evento\nContactCreated"]
        S1 --> S2 --> S3 --> S4
    end

    subgraph REPO["Repository Layer"]
        R1["ContactRepository.Save(contact)"]
        R2[("PostgreSQL\nINSERT INTO contacts...")]
        R1 --> R2
    end

    subgraph ASYNC["Async Workers"]
        E1["Event Bus\n(buffered channel)"]
        W1["ActivityLog Worker\n→ MongoDB"]
        W2["Email Worker\n→ SMTP welcome"]
        E1 --> W1
        E1 --> W2
    end

    H2 -->|"CreateContactDTO válido"| S1
    S4 -->|"ContactCreated event"| E1
    S3 -->|"Contact entity"| R1
    R2 -->|"contact com ID + timestamps"| S1
    S1 -->|"ContactResponse"| H1
    H1 -->|"201 Created"| RESP["HTTP Response\n{id, name, email, created_at}"]
```

### 8.3 Cache-Aside Pattern (Módulo 17)

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant Cache as Redis Cache
    participant DB as PostgreSQL

    Client->>API: GET /contacts/123

    API->>Cache: GET contact:123
    alt Cache HIT
        Cache-->>API: contact JSON
        API-->>Client: 200 OK (< 5ms)
        Note over Client,Cache: Servido do cache
    else Cache MISS
        Cache-->>API: nil
        API->>DB: SELECT * FROM contacts WHERE id = 123
        DB-->>API: contact record
        API->>Cache: SET contact:123 {data} EX 300
        Note over Cache: TTL: 5 minutos
        API-->>Client: 200 OK (~50ms)
    end

    Note over Client,DB: Invalidação do cache

    Client->>API: PUT /contacts/123 {name: "Novo Nome"}
    API->>DB: UPDATE contacts SET name = ? WHERE id = 123
    DB-->>API: updated
    API->>Cache: DEL contact:123
    Note over Cache: Cache invalidado — próximo GET vai ao DB
    API-->>Client: 200 OK
```

### 8.4 Jobs Assíncronos com Goroutines (Módulo 17)

```mermaid
flowchart TD
    subgraph MAIN["Main Application"]
        API["API Server"]
        BUS["Event Bus\nbuffered channel\ncap: 1000"]
    end

    subgraph WORKERS["Worker Pool — goroutines"]
        W1["📧 Email Worker\ngoroutine 1\nSMTP notifications"]
        W2["📝 Activity Logger\ngoroutine 2\nMongoDB logs"]
        W3["📅 Follow-up Scheduler\ngoroutine 3\nDue date reminders"]
        W4["📊 Analytics Worker\ngoroutine 4\nMetrics aggregation"]
    end

    EVENTS["Eventos\n─────────\nContactCreated\nDealWon\nDealLost\nTaskOverdue\nLeadConverted"]

    API -->|"emit event"| BUS
    BUS --> W1
    BUS --> W2
    BUS --> W3
    BUS --> W4
    EVENTS -.->|"define tipos"| BUS

    style WORKERS fill:#fef9c3
    style BUS fill:#dcfce7
```

---

## 9. Design Patterns Aplicados

```mermaid
flowchart TD
    ROOT(["🎨 Design Patterns no GoRM"])

    ROOT --> CR["🏗️ Criacionais"]
    ROOT --> ES["🔗 Estruturais"]
    ROOT --> CO["⚡ Comportamentais"]

    CR --> CR1["Factory\n─ ContactFactory\n─ LeadFactory\n─ UserFactory (testes)"]
    CR --> CR2["Builder\n─ QueryBuilder (filtros)\n─ EmailBuilder (templates)"]
    CR --> CR3["Singleton\n─ DB Connection pool\n─ Redis client"]

    ES --> ES1["Repository\n─ ContactRepository interface\n─ LeadRepository interface\n─ Impl PG + Mock"]
    ES --> ES2["Adapter\n─ EmailAdapter\n─ SMTPAdapter\n─ SendGridAdapter"]
    ES --> ES3["Decorator\n─ LoggingMiddleware\n─ AuthMiddleware\n─ RateLimitMiddleware"]
    ES --> ES4["Facade\n─ CRMFacade"]

    CO --> CO1["Observer\n─ Event Bus\n─ ActivityLog listener\n─ Email listener"]
    CO --> CO2["Strategy\n─ Pipeline Stage\n─ Pricing\n─ Export CSV/JSON"]
    CO --> CO3["Command\n─ Undo tasks\n─ Bulk operations"]
    CO --> CO4["Chain of Resp.\n─ Middleware pipeline\n─ Validation chain"]
    CO --> CO5["Template Method\n─ Base Report\n─ SalesReport\n─ ActivityReport"]

    style ROOT fill:#1e40af,color:#fff
    style CR fill:#065f46,color:#fff
    style ES fill:#92400e,color:#fff
    style CO fill:#4c1d95,color:#fff
```

### 9.1 Repository Pattern em Go

```mermaid
classDiagram
    class ContactRepository {
        <<interface>>
        +FindByID(id UUID) Contact, error
        +FindAll(filters ContactFilters) []Contact, error
        +Save(contact Contact) Contact, error
        +Update(contact Contact) Contact, error
        +Delete(id UUID) error
    }

    class PostgreSQLContactRepository {
        -db *gorm.DB
        +FindByID(id UUID) Contact, error
        +FindAll(filters ContactFilters) []Contact, error
        +Save(contact Contact) Contact, error
        +Update(contact Contact) Contact, error
        +Delete(id UUID) error
    }

    class MockContactRepository {
        -contacts map[UUID]Contact
        +FindByID(id UUID) Contact, error
        +FindAll(filters ContactFilters) []Contact, error
        +Save(contact Contact) Contact, error
        +Update(contact Contact) Contact, error
        +Delete(id UUID) error
    }

    class ContactService {
        -repo ContactRepository
        -eventBus EventBus
        +CreateContact(dto CreateContactDTO) ContactResponse, error
        +GetContact(id UUID) ContactResponse, error
        +UpdateContact(id UUID, dto UpdateContactDTO) ContactResponse, error
        +DeleteContact(id UUID) error
    }

    ContactRepository <|.. PostgreSQLContactRepository : implements
    ContactRepository <|.. MockContactRepository : implements
    ContactService --> ContactRepository : depends on interface
```

---

## 10. Pirâmide de Testes

```mermaid
flowchart TB
    subgraph PYRAMID["Pirâmide de Testes — GoRM"]
        E2E["🔺 E2E / Funcionais — 5%\n─────────────────────────\nTestcontainers + HTTP client real\nFluxo completo: login → criar deal → fechar negócio\nLento · Frágil · Alto valor"]

        INTEG["🔶 Integração — 15%\n─────────────────────────\nTestcontainers (PG + Redis + Mongo reais)\nRepository layer — queries SQL reais\nService layer — com DB real\nMédio · Fiável · Bom custo-benefício"]

        UNIT["🟩 Unitários — 80%\n─────────────────────────\nMocks para dependências externas\nService logic e regras de negócio\nValidações e mapeamentos\nRápido · Isolado · Barato"]
    end

    UNIT --> INTEG --> E2E

    style UNIT fill:#22c55e,color:#fff,stroke:#16a34a
    style INTEG fill:#f59e0b,color:#fff,stroke:#d97706
    style E2E fill:#ef4444,color:#fff,stroke:#dc2626
```

### 10.1 Estratégia de Testes por Módulo

| Módulo | Tipo de testes introduzidos | Ferramenta |
|--------|----------------------------|------------|
| M03+ | Unitários básicos | `testing` nativo + `testify` |
| M06 | Testes de auth (JWT) | `testify/mock` |
| M14 | Integração com DB | `testcontainers-go` |
| M14 | E2E via HTTP | `net/http/httptest` |
| M17 | Benchmarks | `testing.B` |

---

## 11. Pipeline CI/CD

```mermaid
flowchart LR
    DEV["👨‍💻 Developer\ngit push origin feature/xyz"] --> PR

    subgraph PR["Pull Request"]
        REVIEW["Code Review\n(obrigatório)"]
        CHECKS["Status Checks\n(obrigatórios)"]
    end

    subgraph CI["⚙️ GitHub Actions — CI"]
        LINT["golangci-lint\nCode quality + style"]
        UNIT_T["go test -race ./...\nUnit tests + race detector"]
        INT_T["Integration tests\ntestcontainers"]
        BUILD["go build ./cmd/api\nBinary compilado"]
        DOCKER_B["docker build\nImagem válida"]
        LINT --> UNIT_T --> INT_T --> BUILD --> DOCKER_B
    end

    PR --> CI

    CI -->|"✅ Todos os checks passaram"| MERGE["Merge to main"]
    CI -->|"❌ Algum check falhou"| BLOCK["❌ Merge bloqueado"]

    subgraph CD["🚀 Deploy Pipeline"]
        REGISTRY["Push Image\nContainer Registry\n(ECR / GCR)"]
        STAGING["Deploy Staging\nambiente de testes"]
        SMOKE["Smoke Tests\nGET /health · auth flow"]
        APPROVAL["Manual Approval\n(opcional para prod)"]
        PROD["Deploy Production\nBlue/Green ou Rolling"]
        REGISTRY --> STAGING --> SMOKE --> APPROVAL --> PROD
    end

    MERGE --> CD
    PROD --> NOTIFY["📬 Notificação\nSlack / Email"]
```

---

## 12. Plano de Módulos

### Módulo 01 — Setup & Estrutura Go

**Branch:** `branch-01-setup`
**Feature no GoRM:** Projeto Go inicializado, `GET /health` funcional
**Duração estimada:** 3 dias

**Conteúdo:**
- Instalar Go 1.22+, configurar workspace
- Standard Go layout — porquê esta estrutura
- `go mod init`, `go.mod`, `go.sum`
- Primeiro handler HTTP com Fiber
- Makefile para comandos comuns

```mermaid
flowchart LR
    A["go mod init\ngithub.com/user/gorm-crm"] --> B["Standard Layout\ncmd/ internal/ pkg/"]
    B --> C["cmd/api/main.go\nFiber setup"]
    C --> D["GET /health\nretorna 200 OK"]
    D --> E["make run\nlocalhost:8080 ✅"]
```

---

### Módulo 02 — Fundamentos Go

**Branch:** `branch-02-go-fundamentos`
**Feature no GoRM:** Domain models definidos, interfaces de repositório criadas
**Duração estimada:** 5 dias

**Conteúdo (80/20 Go para quem vem de outra linguagem):**
- Structs, interfaces, embedding vs herança
- Error handling idiomático (`error` como valor, não exceção)
- Goroutines e channels — introdução prática
- `defer`, `panic`, `recover`
- Ponteiros vs valores — quando usar cada um
- Generics — introdução com casos reais

```mermaid
flowchart TD
    subgraph MODELS["Domain Models criados"]
        C["Contact struct"]
        L["Lead struct"]
        D["Deal struct"]
        T["Task struct"]
        U["User struct"]
    end
    subgraph INTERFACES["Interfaces definidas"]
        CR["ContactRepository\ninterface"]
        LR["LeadRepository\ninterface"]
        DR["DealRepository\ninterface"]
    end
    MODELS --> INTERFACES
```

---

### Módulo 03 — SQL & PostgreSQL

**Branch:** `branch-03-sql`
**Feature no GoRM:** CRUD completo de Contactos persistido em PostgreSQL
**Duração estimada:** 5 dias

**Conteúdo:**
- PostgreSQL setup com Docker
- GORM — models, migrations, associations
- `golang-migrate` para versioning de schema
- CRUD de Contacts com todos os campos
- Transações ACID — quando e como usar
- Query com filtros, ordenação e paginação
- Índices — porque importam para performance

```mermaid
sequenceDiagram
    participant Handler
    participant Service
    participant Repository
    participant GORM
    participant PG as PostgreSQL

    Handler->>Service: CreateContact(dto)
    Service->>Service: Validar + mapear para entity
    Service->>Repository: Save(contact)
    Repository->>GORM: db.Create(&contact)
    GORM->>PG: INSERT INTO contacts (name, email, ...) VALUES (...)
    PG-->>GORM: {id: uuid, created_at: timestamp}
    GORM-->>Repository: contact preenchido
    Repository-->>Service: Contact
    Service-->>Handler: ContactResponse DTO
    Handler-->>Client: 201 Created {contact}
```

---

### Módulo 07 — Arquitetura MVC em Camadas

**Branch:** `branch-07-mvc-layers`
**Feature no GoRM:** Todos os domínios separados em Handler/Service/Repository
**Duração estimada:** 4 dias

**Conteúdo:**
- Separação de responsabilidades — porquê interessa
- Handler layer — só HTTP, nada de negócio
- Service layer — regras de negócio, sem HTTP
- Repository layer — só acesso a dados, sem negócio
- DTOs vs Domain Models — a fronteira
- Dependency Injection manual (constructors)

```mermaid
flowchart TD
    subgraph HANDLER["Handler — HTTP only"]
        H["• Bind request\n• Validate input\n• Call service\n• Map to response\n• Return HTTP status"]
    end
    subgraph SERVICE["Service — Business only"]
        S["• Orchestrar operações\n• Aplicar regras de negócio\n• Chamar repositórios\n• Emitir eventos\n• Nada de HTTP aqui"]
    end
    subgraph REPO["Repository — Data only"]
        R["• Queries ao DB\n• Mapeamento ORM\n• Nada de negócio\n• Nada de HTTP\n• Retorna domain models"]
    end
    HANDLER -->|"DTO (input validado)"| SERVICE
    SERVICE -->|"Domain entity"| REPO
    REPO -->|"Domain entity"| SERVICE
    SERVICE -->|"Domain entity"| HANDLER
```

---

### Módulo 12 — SOLID em Go

**Branch:** `branch-12-solid`
**Feature no GoRM:** Refactor da codebase aplicando os 5 princípios
**Duração estimada:** 5 dias

```mermaid
flowchart LR
    ROOT(["🏛️ SOLID no GoRM"])

    ROOT --> S["S — Single Responsibility"]
    S --> S1["ContactService só gere contactos"]
    S --> S2["EmailService só envia emails"]
    S --> S3["LogService só escreve logs"]
    S --> S4["Cada struct tem um motivo para mudar"]

    ROOT --> O["O — Open/Closed"]
    O --> O1["Adicionar novo stage no pipeline"]
    O --> O2["sem modificar código existente"]
    O --> O3["Strategy interface para stages"]

    ROOT --> L["L — Liskov Substitution"]
    L --> L1["MockContactRepository"]
    L --> L2["substitui PostgreSQLContactRepo"]
    L --> L3["nos testes sem quebrar comportamento"]

    ROOT --> I["I — Interface Segregation"]
    I --> I1["ContactReader interface"]
    I --> I2["ContactWriter interface"]
    I --> I3["em vez de uma interface gigante"]

    ROOT --> D["D — Dependency Inversion"]
    D --> D1["Services dependem de interfaces"]
    D --> D2["não de implementações concretas"]
    D --> D3["Injectado via constructor"]

    style ROOT fill:#1e40af,color:#fff
    style S fill:#065f46,color:#fff
    style O fill:#92400e,color:#fff
    style L fill:#1e3a8a,color:#fff
    style I fill:#4c1d95,color:#fff
    style D fill:#831843,color:#fff
```

---

### Módulo 14 — Testes Automatizados

**Branch:** `branch-14-testes`
**Feature no GoRM:** Cobertura completa da pirâmide de testes
**Duração estimada:** 5 dias

**Conteúdo:**
- Testes unitários com `testify` e mocks manuais
- Testes de integração com `testcontainers-go`
- Testes E2E com `httptest`
- Table-driven tests (idioma Go)
- Test coverage — `go test -cover`
- Race detector — `go test -race`
- Benchmarks — `go test -bench`

---

### Módulo 15 — Design Patterns

**Branch:** `branch-15-patterns`
**Feature no GoRM:** 10+ patterns aplicados a casos reais
**Duração estimada:** 6 dias

**Patterns e onde aparecem no GoRM:**

| Pattern | Categoria | Onde no GoRM |
|---------|-----------|--------------|
| Repository | Estrutural | `ContactRepository`, `DealRepository` |
| Factory | Criacional | `ContactFactory`, `UserFactory` (testes) |
| Builder | Criacional | `QueryBuilder` para filtros/paginação |
| Observer | Comportamental | `EventBus` para `ActivityLog` e emails |
| Strategy | Comportamental | Pipeline stages, export formats |
| Decorator | Estrutural | `LoggingMiddleware`, `AuthMiddleware` |
| Singleton | Criacional | DB connection pool, Redis client |
| Command | Comportamental | Undo em tasks, bulk operations |
| Facade | Estrutural | `CRMService` agrega sub-serviços |
| Adapter | Estrutural | `EmailAdapter` (SMTP vs SendGrid) |

---

## 13. Roadmap de Implementação

```mermaid
gantt
    title GoRM — Plano de Construção (estimativa self-paced)
    dateFormat  YYYY-MM-DD
    section Nível Júnior
    M01 Setup e Estrutura        :m01, 2026-06-01, 3d
    M02 Fundamentos Go           :m02, after m01, 5d
    M03 SQL e PostgreSQL         :m03, after m02, 5d
    M04 Git Workflow             :m04, after m03, 2d
    M05 REST API                 :m05, after m04, 5d
    M06 Autenticação             :m06, after m05, 4d
    M07 Arquitetura MVC          :m07, after m06, 4d
    M08 Docker                   :m08, after m07, 3d
    section Nível Pleno
    M09 NoSQL MongoDB            :m09, after m08, 4d
    M10 Clean Code               :m10, after m09, 4d
    M11 OOP Avançado             :m11, after m10, 4d
    M12 SOLID                    :m12, after m11, 5d
    M13 Object Calisthenics      :m13, after m12, 3d
    M14 Testes Automatizados     :m14, after m13, 5d
    M15 Design Patterns          :m15, after m14, 6d
    section Nível Sénior
    M16 Refactoring              :m16, after m15, 4d
    M17 Performance e Cache      :m17, after m16, 5d
    M18 Cloud e CI/CD            :m18, after m17, 5d
```

**Total estimado:** ~80 horas de estudo autodirigido (~31 dias úteis)

---

## 14. Decisões de Design (ADRs)

Ver pasta [`docs/adr/`](adr/) para os registos completos.

| ADR | Decisão | Estado |
|-----|---------|--------|
| [ADR-001](adr/001-http-framework.md) | Usar Fiber como framework HTTP | ✅ Aceite |
| [ADR-002](adr/002-orm-strategy.md) | GORM nos módulos iniciais, sqlx no de performance | ✅ Aceite |
| [ADR-003](adr/003-branch-strategy.md) | Branches cumulativas (cada uma parte da anterior) | ✅ Aceite |
| [ADR-004](adr/004-dependency-injection.md) | Manual DI nos primeiros módulos, Wire introduzido mais tarde | ✅ Aceite |
| [ADR-005](adr/005-error-handling.md) | Errors como valores (idioma Go), sem panic em produção | ✅ Aceite |
| [ADR-006](adr/006-testing-strategy.md) | Testcontainers para integração, mocks manuais para unitários | ✅ Aceite |

---

## 15. Checklist por Módulo

Cada branch/módulo deve ter obrigatoriamente:

```
[ ] README.md com:
    [ ] Objetivo claro do módulo
    [ ] Lista de conceitos abordados
    [ ] O que vais construir (feature no GoRM)
    [ ] Diagrama Mermaid de contexto
    [ ] Pré-requisitos (branch anterior)
    [ ] Recursos e referências

[ ] Código Go:
    [ ] Funcional e testável
    [ ] Segue a estrutura de pastas definida
    [ ] Incremento sobre a branch anterior
    [ ] go vet e golangci-lint sem erros

[ ] Testes (a partir do M03):
    [ ] Mínimo unitários para a nova lógica
    [ ] Integração quando há acesso a DB

[ ] Documentação:
    [ ] CHALLENGE.md com exercício prático
    [ ] ADR se houve decisão de design relevante
    [ ] Comentários no código apenas onde necessário

[ ] Git:
    [ ] Conventional commits ao longo do módulo
    [ ] git tag vX.X no final
    [ ] PR documentado para merge em main
```

---

## Resumo Executivo

| Dimensão | Detalhe |
|----------|---------|
| **Nome** | GoRM — Um CRM construído em Go |
| **Total de módulos** | 18 módulos em 3 níveis |
| **Duração estimada** | ~80 horas (self-paced) |
| **Estrutura** | 1 branch Git por módulo, cumulativas |
| **App produzida** | CRM completo: contactos, leads, deals, tasks, auth, cache, deploy |
| **Stack** | Go + Fiber + PostgreSQL + MongoDB + Redis + Docker + GitHub Actions |
| **Nível de entrada** | Programador experiente noutra linguagem |
| **Nível de saída** | Backend Sénior em Go com visão de arquitetura |
| **Reutilizabilidade** | Repositório auto-suficiente para ensinar outros |

---

> *"O melhor código é o código que tu próprio construíste, entendes e consegues explicar."*
