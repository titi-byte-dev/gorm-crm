# CRM OS — Frontend Plan

> Workspace unificado que combina Notion + ClickUp + n8n sobre o backend GoRM CRM.

---

## Visão

Um único workspace onde o utilizador muda de **perspectiva** sobre os mesmos dados — sem alternar entre ferramentas.

```
┌─────────────────────────────────────────────────────────────────┐
│  Object Lens       →  Notion     (páginas ricas, relações)      │
│  Execution Lens    →  ClickUp    (kanban, tasks, timelines)     │
│  Automation Lens   →  n8n/Make   (flows visuais, triggers)      │
└─────────────────────────────────────────────────────────────────┘
```

---

## Tech Stack

| Camada | Tecnologia | Razão |
|--------|-----------|-------|
| Framework | Next.js 15 (App Router) | SSR/SSG, routing, server actions |
| Linguagem | TypeScript strict | Segurança de tipos end-to-end |
| UI base | shadcn/ui + Tailwind CSS | Composable, sem opinião de design |
| State server | TanStack Query v5 | Cache, sync, optimistic updates |
| State client | Zustand | UI state local (sidebars, modals) |
| Forms | React Hook Form + Zod | Validação partilhada com backend |
| Rich text | Tiptap | Extensível, mantido, suporta blocos |
| Flow builder | React Flow | Base para o Automation Lens |
| Drag & Drop | @dnd-kit | Acessível, keyboard-navigable |
| Charts | Recharts | Composable, leve |
| Testes | Vitest + Playwright | Unit + E2E |

---

## Estrutura de Pastas

```
app/
├── (auth)/
│   ├── login/
│   └── register/
└── (workspace)/
    ├── layout.tsx              # Sidebar + TopBar + providers
    ├── dashboard/
    ├── contacts/
    │   ├── page.tsx            # Lista (Execution Lens default)
    │   └── [id]/
    │       └── page.tsx        # Object Lens — página rica
    ├── leads/
    ├── deals/
    ├── tasks/
    ├── automations/
    └── settings/

components/
├── ui/                         # shadcn primitives
├── navigation/
│   ├── Sidebar.tsx
│   ├── TopBar.tsx
│   └── CommandPalette.tsx      # Cmd+K
├── lenses/
│   ├── ObjectLens.tsx          # Tiptap + relações inline
│   ├── ExecutionLens.tsx       # Kanban / Timeline
│   └── AutomationLens.tsx      # React Flow canvas
├── entity/
│   ├── ContactCard.tsx
│   ├── DealCard.tsx
│   └── ActivityFeed.tsx
└── blocks/                     # Notion-like block types

lib/
├── api.ts                      # fetch wrapper tipado
├── auth.ts                     # JWT handling
└── hooks/
    ├── useContacts.ts
    ├── useDeals.ts
    └── useOrganization.ts
```

---

## Módulos e Funcionalidades

### 1. Object Lens — Página de Entidade Rica

Inspirado no Notion. Cada contacto/deal/lead tem uma página com:

- **Header**: nome, avatar, badges de estado
- **Properties panel** (direita): campos editáveis inline
- **Rich document**: blocos Tiptap — texto, headings, checklists, @mentions
- **Related objects**: deals, leads, tasks associados em mini-cards
- **Activity feed**: timeline unificada de eventos

```tsx
// Exemplo de uso
<EntityPage entity={contact}>
  <PropertiesPanel />
  <RichDocument content={contact.notes} />
  <RelatedDeals contactId={contact.id} />
  <ActivityFeed entityId={contact.id} />
</EntityPage>
```

### 2. Execution Lens — Pipeline & Tasks

Inspirado no ClickUp:

- **Kanban** de Deals por stage (drag com @dnd-kit)
- **Lista** de Tasks com filtros por prioridade/estado
- **Split view**: Object Lens | Execution Lens lado a lado
- Manager vê toda a organização; seller vê apenas os seus

### 3. Automation Lens — Flow Builder

Inspirado no n8n/Make. Mapeia diretamente o `events.Bus` do backend:

```
Trigger (event)  →  Condition  →  Action
─────────────────────────────────────────
lead.converted   →  always     →  create task "Acompanhamento"
deal.won         →  value>5000 →  send email + notify Slack
task.overdue     →  always     →  escalate to manager
```

Componentes:
- **FlowCanvas** (React Flow): nodes arrastáveis, edges com labels
- **TriggerNode**: seleciona evento do bus (`contact.created`, `deal.won`, etc.)
- **ActionNode**: cria task, envia notificação, atualiza campo
- **ConditionNode**: filtro por campo/valor
- **FlowSidebar**: lista de flows ativos + histórico de execuções

### 4. Command Palette (Cmd+K)

Acessível globalmente:
- Navegar para qualquer entidade
- Criar contacto/lead/deal/task em linha
- Mudar de lens
- Executar ações rápidas ("Marcar deal como ganho")

---

## API Integration

O frontend comunica com o backend GoRM via REST. Todos os pedidos incluem `Authorization: Bearer <token>`.

```typescript
// lib/api.ts
const api = {
  contacts: {
    list: (filters) => fetch('/api/v1/contacts?' + new URLSearchParams(filters)),
    get: (id) => fetch(`/api/v1/contacts/${id}`),
    create: (dto) => fetch('/api/v1/contacts', { method: 'POST', body: JSON.stringify(dto) }),
    update: (id, dto) => fetch(`/api/v1/contacts/${id}`, { method: 'PUT', body: JSON.stringify(dto) }),
    delete: (id) => fetch(`/api/v1/contacts/${id}`, { method: 'DELETE' }),
  },
  // ... leads, deals, tasks, auth
}
```

Multi-tenancy é transparente — o `org_id` viaja no JWT, o backend filtra automaticamente.

---

## Role-Based UI

| Role | Vê | Pode criar/editar |
|------|-----|-------------------|
| `seller` | Só os seus dados | Sim, nos seus |
| `manager` | Toda a organização | Sim, em todos |
| `admin` | Toda a organização | Sim + configurações |

A sidebar e os filtros adaptam-se ao role extraído do JWT.

---

## Acessibilidade (WCAG 2.2 AA)

- Navegação completa por teclado (sidebar, kanban, flow builder)
- Drag & drop com alternativa por botões (WCAG 2.5.7)
- Focus indicators visíveis com contraste 3:1 mínimo
- `prefers-reduced-motion` respeita todas as animações
- Targets mínimos 24×24px (WCAG 2.5.8)
- Testes automáticos com axe-core no CI

---

## Roadmap

### Fase 1 — MVP (semanas 1–6)
- [ ] Auth (login, register, JWT refresh)
- [ ] Sidebar + TopBar + Command Palette
- [ ] Contacts — lista + Object Lens básico
- [ ] Deals — Kanban com drag
- [ ] Tasks — lista com filtros
- [ ] Design System: tokens, componentes base

### Fase 2 — v0.8 (semanas 7–12)
- [ ] Rich document (Tiptap) em entidades
- [ ] Activity Feed unificada
- [ ] Split view Object | Execution
- [ ] Leads pipeline
- [ ] Role-based UI completo

### Fase 3 — v1.0 (semanas 13–20)
- [ ] Automation Lens (Flow Builder básico)
- [ ] Triggers mapeados ao events.Bus
- [ ] Relatórios e dashboards
- [ ] Internacionalização (PT/EN)
- [ ] Testes E2E com Playwright

### Fase 4 — v2.0 (Agent Mode)
- [ ] `AgentButton` por entidade (contact, deal, lead)
- [ ] `AgentPanel` com seletor de tipo + modo suggest/auto
- [ ] `AgentActionCard` — aprovar/rejeitar cada ação sugerida
- [ ] `AgentExecutionLog` — histórico de runs por entidade
- [ ] SSE stream para feedback em tempo real durante execução
- [ ] Client Portal (acesso limitado do cliente)
- [ ] Mobile (React Native ou PWA)
- [ ] Integrações externas (email, Slack, calendário)

---

## 4. Agent Mode — Camada de Inteligência

> "Um agente que age sobre os dados em vez de apenas os mostrar."

O Agent Mode é a quarta camada do CRM OS. Não é uma funcionalidade separada — é um botão **Ativar Agente** disponível em cada entidade (contacto, deal, lead).

### Arquitetura

```
Frontend                         Backend (Go)
────────────────────────────     ─────────────────────────────
AgentButton (entity page)   →   POST /api/v1/agents/run
AgentPanel (Sheet/Drawer)   ←   AgentRun { actions, summary }
AgentActionCard             →   POST /api/v1/agents/runs/:id/approve
AgentExecutionLog           ←   GET  /api/v1/agents/runs?entity_type=&entity_id=
```

### Tipos de Agente

| Tipo | Trigger | O que faz |
|------|---------|-----------|
| `follow_up` | Botão no ContactPage | Analisa contacto + tasks, propõe acompanhamento |
| `deal_closer` | Botão no DealPage | Analisa deal estagnado, propõe próximos passos |
| `task_router` | Botão no TaskBoard | Prioriza e redistribui tasks pendentes |
| `summarize` | Botão em qualquer entidade | Resume histórico em linguagem natural |

### Modos de Execução

```
suggest (padrão)          auto (manager+ apenas)
──────────────────         ──────────────────────
LLM decide ações          LLM decide ações
↓                          ↓
Devolve como sugestões    Executa imediatamente
↓                          ↓
Utilizador aprova/rejeita  Log no ActivityFeed
```

### Componentes Frontend

```tsx
// Botão que abre o painel — aparece no header de cada entidade
<AgentButton entityType="contact" entityId={contact.id} />

// Painel lateral (Sheet)
<AgentPanel>
  <AgentTypeSelector />        {/* follow_up | deal_closer | ... */}
  <ModeToggle />               {/* suggest | auto (manager apenas) */}
  <AgentThinking />            {/* skeleton durante chamada LLM */}
  {run.actions.map(action =>
    <AgentActionCard
      action={action}
      onApprove={approveAction}
      onReject={rejectAction}
    />
  )}
  <AgentSummary text={run.summary} />
</AgentPanel>

// Histórico de execuções
<AgentExecutionLog entityType="contact" entityId={contact.id} />
```

### Estrutura de pastas

```
components/
└── agents/
    ├── AgentButton.tsx
    ├── AgentPanel.tsx
    ├── AgentTypeSelector.tsx
    ├── AgentActionCard.tsx
    ├── AgentExecutionLog.tsx
    └── AgentThinking.tsx      # skeleton de loading

lib/
└── hooks/
    ├── useAgentRun.ts         # POST /agents/run + polling/SSE
    └── useAgentHistory.ts     # GET /agents/runs?entity_type=&entity_id=

stores/
└── agentStore.ts              # Zustand: activeRun, pendingActions
```

### Role-Based Agent UI

| Role | ModeSuggest | ModeAuto | Ver histórico |
|------|-------------|----------|---------------|
| `seller` | ✅ | ❌ | Só os seus runs |
| `manager` | ✅ | ✅ | Toda a org |
| `admin` | ✅ | ✅ | Toda a org |

### Fallback sem API Key

Se `ANTHROPIC_API_KEY` não estiver configurada no backend, o agente corre em **modo regra** (heurísticas simples):
- Contacto sem atualização há ≥7 dias → cria task follow-up urgente
- Task em atraso → escala ao manager
- Sem condição → devolve sumário vazio

O frontend não precisa saber — o comportamento é transparente.

---

## Integração com o Backend

O backend GoRM expõe todos os endpoints necessários:

| Frontend feature | Backend endpoint |
|-----------------|-----------------|
| Login | `POST /api/v1/auth/login` |
| Register (cria org) | `POST /api/v1/auth/register` |
| Lista de contactos | `GET /api/v1/contacts` |
| Kanban de deals | `GET /api/v1/deals?stage=...` |
| Mover deal | `PATCH /api/v1/deals/:id/stage` |
| Tasks do utilizador | `GET /api/v1/tasks` |
| Tasks em atraso | `GET /api/v1/tasks/overdue` |
| Automation triggers | `events.Bus` (futuro: WebSocket ou SSE) |
| Ativar agente | `POST /api/v1/agents/run` |
| Aprovar ações | `POST /api/v1/agents/runs/:id/approve` |
| Histórico de runs | `GET /api/v1/agents/runs?entity_type=&entity_id=` |

O `events.Bus` do backend é a ponte natural para o Automation Lens — cada evento publicado (`deal.won`, `lead.converted`, `task.completed`, `agent.run.completed`) pode disparar um flow visual.
