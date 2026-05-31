# 🎨 GoRM CRM — Design System

> Tokens de design do produto GoRM CRM.
> Framework-agnostic — definidos uma vez, consumíveis por qualquer frontend.

---

## Identidade Visual

**GoRM** é um CRM profissional construído em Go. A identidade visual reflete os valores do produto:

| Valor | Expressão visual |
|-------|-----------------|
| **Técnico mas acessível** | Tipografia clara (Inter), espaçamentos generosos |
| **Confiança** | Azul primário — associado a profissionalismo e fiabilidade |
| **Energia** | Violeta secundário — diferenciador, não genérico |
| **Clareza** | Hierarquia visual forte, nunca decoração por decoração |

---

## Paleta de Cores

### Brand

```mermaid
flowchart LR
    P["⬛ Primary\n#0F62FE\nAzul IBM"]
    S["⬛ Secondary\n#6929C4\nVioleta"]
    A["⬛ Accent\n#00B4D8\nCiano"]

    style P fill:#0F62FE,color:#fff,stroke:#0F62FE
    style S fill:#6929C4,color:#fff,stroke:#6929C4
    style A fill:#00B4D8,color:#fff,stroke:#00B4D8
```

### Pipeline de Vendas — Cores por Estado

Cada estado do pipeline tem uma cor semântica que comunica progresso e urgência:

```mermaid
flowchart LR
    N["New\n#6B7280"]
    C["Contacted\n#3B82F6"]
    Q["Qualified\n#F59E0B"]
    P["Proposal\n#8B5CF6"]
    NG["Negotiation\n#EC4899"]
    W["Won ✅\n#10B981"]
    L["Lost ❌\n#EF4444"]

    N --> C --> Q --> P --> NG --> W
    NG --> L
    Q --> L

    style N  fill:#6B7280,color:#fff,stroke:#6B7280
    style C  fill:#3B82F6,color:#fff,stroke:#3B82F6
    style Q  fill:#F59E0B,color:#fff,stroke:#F59E0B
    style P  fill:#8B5CF6,color:#fff,stroke:#8B5CF6
    style NG fill:#EC4899,color:#fff,stroke:#EC4899
    style W  fill:#10B981,color:#fff,stroke:#10B981
    style L  fill:#EF4444,color:#fff,stroke:#EF4444
```

### Roles de Utilizador

```mermaid
flowchart LR
    AD["👑 Admin\n#6929C4"]
    MG["👔 Manager\n#0F62FE"]
    SL["🧑‍💼 Seller\n#00B4D8"]

    style AD fill:#6929C4,color:#fff,stroke:#6929C4
    style MG fill:#0F62FE,color:#fff,stroke:#0F62FE
    style SL fill:#00B4D8,color:#fff,stroke:#00B4D8
```

### Prioridade de Tarefas

```mermaid
flowchart LR
    LW["Low\n#9CA3AF"]
    MD["Medium\n#3B82F6"]
    HG["High\n#F59E0B"]
    UG["Urgent\n#EF4444"]

    style LW fill:#9CA3AF,color:#fff,stroke:#9CA3AF
    style MD fill:#3B82F6,color:#fff,stroke:#3B82F6
    style HG fill:#F59E0B,color:#fff,stroke:#F59E0B
    style UG fill:#EF4444,color:#fff,stroke:#EF4444
```

---

## Tipografia

| Role | Font | Tamanho | Peso |
|------|------|---------|------|
| Headings | Inter | 24–36px | 700 |
| Body | Inter | 16px | 400 |
| Labels / Meta | Inter | 12–14px | 500 |
| Código | JetBrains Mono | 14px | 400 |

---

## Escala de Espaçamento

Base unit: **4px**

| Token | Valor | Uso típico |
|-------|-------|-----------|
| `spacing.1` | 4px | Padding interno de ícone |
| `spacing.2` | 8px | Gap entre elementos inline |
| `spacing.4` | 16px | Padding de card, espaço entre campos |
| `spacing.6` | 24px | Padding de secção |
| `spacing.8` | 32px | Margem entre blocos |
| `spacing.16` | 64px | Padding de página |

---

## Border Radius

| Token | Valor | Uso |
|-------|-------|-----|
| `radius.sm` | 4px | Badges, tags |
| `radius.base` | 6px | Inputs, botões |
| `radius.md` | 8px | Cards |
| `radius.lg` | 12px | Modals, drawers |
| `radius.full` | 9999px | Pills, avatares |

---

## Anatomia de um Card de Deal

```mermaid
flowchart TD
    subgraph CARD["Deal Card — exemplo de aplicação dos tokens"]
        HEADER["Header\nfont: Inter Bold 16px\nbg: color.background.subtle"]
        STAGE["Stage Badge\nbg: color.pipeline.negotiation (#EC4899)\nradius: radius.full\nfont: 12px semibold"]
        VALUE["€ 12.500\nfont: Inter Bold 24px\ncolor: color.brand.primary"]
        OWNER["👤 João Silva\nfont: 14px\ncolor: color.neutral.500"]
        DUE["⏰ Vence em 3 dias\ncolor: color.semantic.warning"]
    end

    style CARD fill:#F9FAFB,stroke:#E5E7EB
    style HEADER fill:#F3F4F6,stroke:none
    style STAGE fill:#EC4899,color:#fff,stroke:none
    style VALUE fill:#F9FAFB,stroke:none
```

---

## Como consumir os tokens

Os tokens estão em [`tokens.json`](tokens.json) num formato standard.

### CSS Custom Properties (quando houver frontend)
```css
:root {
  --color-brand-primary:   #0F62FE;
  --color-pipeline-won:    #10B981;
  --color-pipeline-lost:   #EF4444;
  --font-sans:             'Inter', system-ui, sans-serif;
  --font-mono:             'JetBrains Mono', monospace;
  --spacing-4:             16px;
  --radius-md:             8px;
}
```

### Tailwind Config (quando houver frontend)
```js
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        brand: { primary: '#0F62FE', secondary: '#6929C4' },
        pipeline: {
          new: '#6B7280', contacted: '#3B82F6',
          won: '#10B981',  lost: '#EF4444',
        }
      }
    }
  }
}
```

### API — Campos de cor nas respostas
O backend devolve o `stage`/`status` como string. O frontend faz o mapeamento local usando os tokens. A API **nunca devolve cores** — separação de responsabilidades.

```json
{ "stage": "negotiation" }
// Frontend mapeia: "negotiation" → color.pipeline.negotiation → #EC4899
```

---

## Mermaid Theme para Diagramas do Curso

Todos os diagramas do curso usam esta paleta para consistência:

```
Júnior  → #22c55e (verde)
Pleno   → #3b82f6 (azul)
Sénior  → #a855f7 (violeta)
Neutro  → #e5e7eb (cinzento claro)
```

---

## Ficheiros

| Ficheiro | Conteúdo |
|----------|----------|
| [`tokens.json`](tokens.json) | Todos os tokens em formato standard |
| [`README.md`](README.md) | Este documento — preview visual |

> Quando o frontend for construído, os tokens são importados diretamente deste ficheiro — zero retrabalho.
