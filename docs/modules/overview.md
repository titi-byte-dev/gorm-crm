# Módulos do Curso — Visão Geral

Cada módulo tem a sua própria branch, README e CHALLENGE. Esta página serve de índice rápido.

---

## Mapa de Progressão

```mermaid
flowchart TD
    START([🚀 START]) --> M01

    subgraph JR["🟢 Nível Júnior — branches 01-08"]
        M01["M01 · Setup"]
        M02["M02 · Fundamentos Go"]
        M03["M03 · SQL & PostgreSQL"]
        M04["M04 · Git Workflow"]
        M05["M05 · REST API"]
        M06["M06 · Auth & JWT"]
        M07["M07 · MVC Layers"]
        M08["M08 · Docker"]
        M01-->M02-->M03-->M04-->M05-->M06-->M07-->M08
    end

    M08 --> J(["🏆 Júnior"])

    subgraph PL["🔵 Nível Pleno — branches 09-15"]
        M09["M09 · NoSQL & MongoDB"]
        M10["M10 · Clean Code"]
        M11["M11 · OOP Avançado"]
        M12["M12 · SOLID"]
        M13["M13 · Object Calisthenics"]
        M14["M14 · Testes"]
        M15["M15 · Design Patterns"]
        M09-->M10-->M11-->M12-->M13-->M14-->M15
    end

    J --> M09

    subgraph SR["🟣 Nível Sénior — branches 16-18"]
        M16["M16 · Refactoring"]
        M17["M17 · Performance & Cache"]
        M18["M18 · Cloud & CI/CD"]
        M16-->M17-->M18
    end

    M15 --> M16
    M18 --> S(["🎓 Sénior"])
```

---

## Índice

### 🟢 Nível Júnior

| # | Branch | Conceitos chave | Feature no GoRM |
|---|--------|-----------------|-----------------|
| 01 | `branch-01-setup` | Go layout, go modules, Makefile | `GET /health` |
| 02 | `branch-02-go-fundamentos` | Structs, interfaces, goroutines, errors | Domain models |
| 03 | `branch-03-sql` | PostgreSQL, GORM, migrations, CRUD | CRUD Contactos |
| 04 | `branch-04-git-workflow` | Branching, conventional commits, PRs | Git workflow |
| 05 | `branch-05-rest-api` | Fiber, REST, middlewares, paginação | API REST completa |
| 06 | `branch-06-auth` | JWT, bcrypt, RBAC, refresh tokens | Auth + Roles |
| 07 | `branch-07-mvc-layers` | Handler/Service/Repository, DTOs, DI | Camadas separadas |
| 08 | `branch-08-docker` | Dockerfile, docker-compose, multi-stage | App containerizada |

### 🔵 Nível Pleno

| # | Branch | Conceitos chave | Feature no GoRM |
|---|--------|-----------------|-----------------|
| 09 | `branch-09-nosql` | MongoDB, document store, activity logs | Histórico de ações |
| 10 | `branch-10-clean-code` | Nomes claros, funções pequenas, sem duplicação | Refactor geral |
| 11 | `branch-11-oop` | Embedding, composição, DRY/KISS/YAGNI | Interfaces avançadas |
| 12 | `branch-12-solid` | S/O/L/I/D em Go com casos reais | Refactor SOLID |
| 13 | `branch-13-calisthenics` | 9 regras Object Calisthenics | Código mais expressivo |
| 14 | `branch-14-testes` | Unitários, integração, E2E, testcontainers | Cobertura completa |
| 15 | `branch-15-patterns` | Repository, Observer, Factory, Strategy... | 10+ patterns |

### 🟣 Nível Sénior

| # | Branch | Conceitos chave | Feature no GoRM |
|---|--------|-----------------|-----------------|
| 16 | `branch-16-refactoring` | Extract method, Replace conditional, Move field | Codebase limpa |
| 17 | `branch-17-performance` | Redis, Cache-Aside, goroutines, benchmarks | Cache + workers |
| 18 | `branch-18-cloud-cicd` | GitHub Actions, Docker registry, deploy | CI/CD completo |

---

## Anatomia de cada branch

```
branch-XX-nome/
├── README.md          ← objetivo, conceitos, diagrama
├── CHALLENGE.md       ← exercício prático do módulo
├── internal/          ← código Go (incremento sobre módulo anterior)
├── docs/              ← ADRs e diagramas do módulo
└── tests/             ← testes do módulo
```

---

## Como navegar

```bash
# Ver todos os módulos disponíveis
git branch -a | grep branch-

# Ir para um módulo
git checkout branch-05-rest-api

# Ver o que mudou neste módulo
git diff branch-04-git-workflow..branch-05-rest-api

# Ver só os ficheiros alterados
git diff --name-only branch-04-git-workflow..branch-05-rest-api

# Ver o histórico de commits do módulo
git log --oneline branch-04-git-workflow..branch-05-rest-api
```
