# GoRM CRM — Guia do Programador

> Stack: Go 1.22+ · Fiber v2 · GORM · PostgreSQL · MongoDB · JWT · events.Bus

---

## Visão Geral da Arquitectura

```
cmd/api/
└── main.go              # wiring: cria dependências, liga rotas, arranca servidor

internal/
├── agent/               # Agent Mode — IA sobre entidades CRM
├── auth/                # JWT, bcrypt, middleware de autenticação
├── activitylog/         # Histórico de atividade (MongoDB)
├── contact/             # Domínio de contactos
├── deal/                # Domínio de deals
├── lead/                # Domínio de leads
├── organization/        # Tenant root — criado no registo
├── task/                # Domínio de tasks
└── shared/
    ├── ctxutil/         # RequestCtx — UserID, TenantID, Role extraídos do JWT
    ├── errors/          # Sentinel errors + handler HTTP global
    ├── events/          # Event Bus em memória (channels Go)
    ├── middleware/       # Logger, CORS
    ├── response/        # Helpers JSON: OK, Created, NoContent, Page
    └── validate/        # Wrapper go-playground/validator

pkg/
├── database/            # Conexão PostgreSQL + MongoDB
├── logger/              # slog estruturado (JSON em prod, texto em dev)
└── version/             # Versão e commit (injetados em build time)

migrations/              # SQL puro — up/down por número de migração
tests/unit/              # Testes unitários (mocks manuais, sem framework)
```

### Padrão por domínio

Cada domínio tem sempre 4 ficheiros:

```
internal/<domínio>/
├── model.go          # struct do domínio + Repository interface + Filters
├── repository_pg.go  # implementação GORM (record struct separada do modelo)
├── service.go        # lógica de negócio + DTOs
└── handler.go        # HTTP handlers Fiber + RegisterRoutes()
```

---

## Multi-tenancy

Cada utilizador pertence a uma **organização** (tenant). O `org_id` viaja no JWT e é injetado pelo middleware em `c.Locals(ctxutil.KeyOrgID, orgID)`.

Todos os handlers extraem o contexto com:

```go
rctx, err := ctxutil.FromFiber(c)
// rctx.UserID   — quem fez o pedido
// rctx.TenantID — organização do utilizador
// rctx.Role     — seller | manager | admin
// rctx.IsManager() — true se manager ou admin
```

Regras de acesso:
- **Seller**: vê apenas os seus dados (`owner_id = uid` ou `assigned_to = uid`)
- **Manager/Admin**: vê toda a organização (`tenant_id = org_id`)

O `checkAccess()` em cada service aplica esta lógica após `FindByID`.

---

## Event Bus

O `events.Bus` é um channel Go em memória com buffer de 500 eventos.

### Publicar
```go
bus.Publish(events.Event{
    Type:    events.TaskCompleted,
    Payload: map[string]string{"id": task.ID.String()},
    UserID:  rctx.UserID.String(),
})
```

### Subscrever (ex: activitylog)
```go
bus.Subscribe(events.DealWon, func(ctx context.Context, e events.Event) {
    // gravar no MongoDB
})
```

### Eventos disponíveis

| Evento | Publicado por |
|--------|--------------|
| `contact.created/updated/deleted` | contact.Service |
| `lead.created/converted/lost` | lead.Service |
| `deal.won/lost` | deal.Service |
| `task.completed/overdue` | task.Service / scheduler |
| `agent.run.completed` | agent.Service |

---

## Agent Mode

### Como funciona

```
POST /api/v1/agents/run
       ↓
agent.Service.Run(rctx, dto)
       ↓
loadEntityContext()  ← carrega contact/deal/lead + tasks do PostgreSQL
       ↓
BuildPrompt()        ← constrói prompt com contexto estruturado
       ↓
LLMClient.Run()      ← chama Anthropic API (ou fallback rule-based)
       ↓
parse ToolCalls      ← cria_task | update_lead_status | add_note | …
       ↓
if ModeAuto → Executor.Execute() → task.Service.Create(rctx, dto)
if ModeSuggest → devolve actions com status "pending_approval"
       ↓
repo.Update(run)     ← persiste AgentRun com actions + summary
       ↓
bus.Publish(AgentRunCompleted)
```

### Configuração

```env
ANTHROPIC_API_KEY=sk-ant-...   # obrigatório para IA; sem ela usa regras
ANTHROPIC_MODEL=claude-haiku-4-5-20251001  # opcional, este é o padrão
```

### Adicionar uma nova tool ao agente

1. Adicionar entrada em `crmTools` em `internal/agent/llm.go`
2. Adicionar `case "nome_da_tool":` no `Executor.Execute()` em `internal/agent/executor.go`
3. Atualizar o `systemPrompt` se necessário

### Fallback sem LLM

Se `ANTHROPIC_API_KEY` não estiver configurada, `NewLLMClient()` devolve `nil` e o `Service.Run()` chama `runRuleBased()` — heurísticas simples que não dependem de API externa.

---

## Migrações

As migrações são SQL puro numeradas sequencialmente. Não existe runner automático — aplicar manualmente ou integrar com `golang-migrate`.

```bash
# Aplicar todas as migrações (exemplo com psql)
for f in migrations/*.up.sql; do psql $DATABASE_URL -f $f; done

# Reverter a última
psql $DATABASE_URL -f migrations/007_create_agent_runs.down.sql
```

| Migração | O que faz |
|----------|-----------|
| 001 | Tabela `users` |
| 002 | Tabela `contacts` |
| 003 | Tabela `leads` |
| 004 | Tabela `deals` |
| 005 | Tabela `tasks` |
| 006 | Tabela `organizations` + colunas `tenant_id` em todas as entidades |
| 007 | Tabela `agent_runs` |

---

## Autenticação

JWT com dois tokens:
- **access_token**: 15 minutos, usado em todos os pedidos
- **refresh_token**: 7 dias, usado apenas em `POST /api/v1/auth/refresh`

Claims:
```go
type Claims struct {
    UserID string    `json:"uid"`
    OrgID  string    `json:"org_id"`
    Role   user.Role `json:"role"`
    jwt.RegisteredClaims
}
```

O middleware `auth.Protected()` valida o token e injeta `userID`, `orgID`, `role` em `c.Locals`.

---

## Testes

```bash
# Todos os testes
go test ./...

# Só unitários (sem base de dados)
go test ./tests/unit/...

# Com verbose
go test -v ./tests/unit/...
```

Os testes unitários usam **mocks manuais** (sem mockery ou gomock). Cada mock implementa a interface do repositório diretamente no ficheiro de teste.

Exemplo de mock:
```go
type mockTaskRepository struct {
    tasks map[uuid.UUID]*task.Task
}

func (m *mockTaskRepository) FindAll(_ uuid.UUID, _ uuid.UUID, _ bool, _ task.Filters) ([]*task.Task, int64, error) {
    return nil, 0, nil
}
```

---

## Variáveis de Ambiente

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `ENV` | `development` | `development` ou `production` |
| `PORT` | `8080` | Porta HTTP |
| `DB_HOST` | `localhost` | Host PostgreSQL |
| `DB_PORT` | `5432` | Porta PostgreSQL |
| `DB_USER` | `postgres` | Utilizador PostgreSQL |
| `DB_PASSWORD` | `postgres` | Password PostgreSQL |
| `DB_NAME` | `gorm_crm` | Nome da base de dados |
| `DB_SSLMODE` | `disable` | SSL mode PostgreSQL |
| `MONGO_URI` | `mongodb://localhost:27017` | URI MongoDB (opcional) |
| `MONGO_DB` | `gorm_crm` | Nome da base de dados MongoDB |
| `JWT_SECRET` | — | **Obrigatório em produção** |
| `ANTHROPIC_API_KEY` | — | Opcional — ativa IA no Agent Mode |
| `ANTHROPIC_MODEL` | `claude-haiku-4-5-20251001` | Modelo Anthropic a usar |

---

## Adicionar um novo domínio

1. Criar `internal/<domínio>/model.go` com a struct + Repository interface + Filters
2. Criar `internal/<domínio>/repository_pg.go` com `record` struct + implementação GORM
3. Criar `internal/<domínio>/service.go` com DTOs e lógica de negócio usando `ctxutil.RequestCtx`
4. Criar `internal/<domínio>/handler.go` com `RegisterRoutes()` e handlers
5. Ligar em `cmd/api/main.go`:
   ```go
   <domínio>.RegisterRoutes(protected, <domínio>.NewService(<domínio>.NewPostgresRepository(db), bus))
   ```
6. Criar migração `migrations/NNN_create_<domínio>.up.sql`

---

## Convenções de código

- Erros de domínio: `fmt.Errorf("contexto: %w", sharederrors.ErrNotFound)`
- Nunca retornar `nil, nil` — sempre um erro ou um valor
- Repository: só CRUD puro — sem lógica de negócio
- Service: toda a lógica de negócio — sem HTTP
- Handler: só parsing HTTP → service → resposta
- Sem comentários que explicam O QUÊ — só o PORQUÊ quando não é óbvio
- Sort columns: sempre via whitelist (`allowedSortColumns`) — nunca interpolar strings SQL
