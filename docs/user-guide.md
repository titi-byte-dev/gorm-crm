# GoRM CRM — Guia do Utilizador

> Versão atual: ver `pkg/version/version.go`

---

## O que é o GoRM CRM?

O GoRM CRM é um sistema de gestão de relacionamento com clientes (CRM) para equipas de vendas B2B. Permite gerir contactos, acompanhar leads, fechar deals e organizar tarefas — tudo num único workspace, com controlo de acesso por organização.

---

## Primeiros passos

### Registo

Ao registar, crias automaticamente uma **organização** (o teu tenant isolado). Nenhuma outra empresa acede aos teus dados.

```http
POST /api/v1/auth/register
{
  "name":     "Ana Silva",
  "email":    "ana@empresa.pt",
  "password": "password123",
  "role":     "admin",
  "org_name": "Empresa Lda."
}
```

Roles disponíveis:

| Role | O que vê | O que pode fazer |
|------|----------|------------------|
| `seller` | Só os seus dados | Criar e editar os seus registos |
| `manager` | Toda a organização | Criar, editar e ver todos os registos |
| `admin` | Toda a organização | Tudo + configurações da organização |

### Login

```http
POST /api/v1/auth/login
{
  "email":    "ana@empresa.pt",
  "password": "password123"
}
```

Resposta: `access_token` (válido 15 min) + `refresh_token` (válido 7 dias).

Todos os pedidos seguintes requerem o header:
```
Authorization: Bearer <access_token>
```

---

## Contactos

Os contactos são as pessoas com quem a tua equipa interage.

### Criar contacto
```http
POST /api/v1/contacts
{
  "name":    "João Ferreira",
  "email":   "joao@cliente.pt",
  "phone":   "+351 912 345 678",
  "company": "Cliente Lda.",
  "notes":   "Interessado no plano Enterprise"
}
```

### Listar contactos
```http
GET /api/v1/contacts?page=1&limit=20&search=João&company=Cliente
```

Parâmetros opcionais: `search`, `company`, `sort_by` (`name`, `company`, `created_at`), `sort_dir` (`asc`, `desc`).

### Atualizar / Eliminar
```http
PUT    /api/v1/contacts/:id
DELETE /api/v1/contacts/:id
```

> **Nota de acesso:** um `seller` só vê e edita os seus próprios contactos. Um `manager` vê toda a organização.

---

## Leads

Um lead representa uma oportunidade de venda em qualquer fase inicial.

### Estados possíveis

```
new → contacted → qualified → lost
 ↘                            ↗
  ─────────── lost ───────────
```

### Criar lead
```http
POST /api/v1/leads
{
  "title":      "Interesse em produto X",
  "value":      5000.00,
  "contact_id": "uuid-do-contacto"
}
```

### Mudar estado
```http
PATCH /api/v1/leads/:id/status
{ "status": "contacted" }
```

Transições inválidas (ex: `qualified → new`) são rejeitadas com erro 422.

---

## Deals

Um deal é uma oportunidade de venda em negociação ativa.

### Pipeline

```
proposal → negotiation → won
        ↘              ↗
         ───── lost ────
```

### Criar deal
```http
POST /api/v1/deals
{
  "title":      "Contrato anual — Empresa X",
  "value":      12000.00,
  "contact_id": "uuid-do-contacto",
  "lead_id":    "uuid-do-lead"  // opcional
}
```

### Mover no pipeline
```http
PATCH /api/v1/deals/:id/stage
{ "stage": "negotiation" }
```

---

## Tasks

As tasks são ações concretas associadas a contactos ou deals.

### Prioridades: `low` | `medium` | `high` | `urgent`
### Estados: `todo` → `in_progress` → `done` | `cancelled`

Uma task no estado `done` ou `cancelled` **não pode ser reaberta**.

### Criar task
```http
POST /api/v1/tasks
{
  "title":       "Enviar proposta comercial",
  "priority":    "high",
  "assigned_to": "uuid-do-utilizador",
  "contact_id":  "uuid-do-contacto",
  "due_date":    "2026-06-15"
}
```

### Tasks em atraso
```http
GET /api/v1/tasks/overdue
```

---

## Agent Mode

O Agent Mode analisa uma entidade (contacto, deal ou lead) e propõe ou executa ações automaticamente.

### Ativar um agente
```http
POST /api/v1/agents/run
{
  "agent_type":  "follow_up",
  "entity_type": "contact",
  "entity_id":   "uuid-do-contacto",
  "mode":        "suggest"
}
```

### Tipos de agente

| Tipo | Para usar quando... |
|------|---------------------|
| `follow_up` | Queres saber o próximo passo com um contacto |
| `deal_closer` | O deal está parado e não sabes como avançar |
| `task_router` | Tens muitas tasks e queres priorização automática |
| `summarize` | Queres um resumo rápido de uma entidade |

### Modos

- **`suggest`** (padrão): o agente propõe ações — tu aprovas ou rejeitas cada uma.
- **`auto`** (manager/admin): o agente executa diretamente. Requer aprovação do teu administrador.

### Aprovar ações sugeridas
```http
POST /api/v1/agents/runs/:run_id/approve
{
  "action_indices": [0, 2]  // índices das ações a executar
}
```

### Ver histórico de runs
```http
GET /api/v1/agents/runs?entity_type=contact&entity_id=uuid
```

> **Sem API Key:** se o administrador não configurou `ANTHROPIC_API_KEY`, o agente funciona em modo regra simples (heurísticas básicas sem IA). O comportamento é transparente — a resposta tem o mesmo formato.

---

## Histórico de Atividade

O CRM regista automaticamente todas as ações importantes no historial de atividade (requer MongoDB configurado).

```http
GET /api/v1/activity/:entity_type/:entity_id
```

---

## Health Check

```http
GET /health
```

Devolve o estado da base de dados PostgreSQL e MongoDB.

---

## Erros comuns

| Código | Significado | O que fazer |
|--------|-------------|-------------|
| 401 | Token expirado ou inválido | Fazer refresh (`POST /api/v1/auth/refresh`) |
| 403 | Sem permissão | Verificar o teu role |
| 404 | Entidade não encontrada | Verificar o ID ou se pertence à tua organização |
| 409 | Email já registado | Usar outro email |
| 422 | Dados inválidos | Ver o campo `details` na resposta de erro |
