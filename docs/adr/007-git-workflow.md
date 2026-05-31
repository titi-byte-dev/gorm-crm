# ADR-007 — Git Workflow

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 04 — Git Workflow

---

## Contexto

O repositório serve dois propósitos em simultâneo: é um curso navegável por branches e é uma app real em desenvolvimento. O workflow Git precisa de suportar ambos.

## Decisão

**Branch-per-module** com merge para `main` via PR documentado.

```
main                    ← estado mais recente e estável
  └── branch-01-setup   ← cada módulo é uma branch cumulativa
  └── branch-02-go-fundamentos
  └── branch-03-sql
  └── ...
  └── branch-18-cloud-cicd
```

Para trabalho pessoal dentro de um módulo (desafios, experiências):

```
branch-03-sql
  └── meu-desafio-m03   ← branch pessoal, nunca entra no main
```

## Conventional Commits

Formato obrigatório: `<tipo>(<escopo>): <descrição>`

| Tipo | Quando usar |
|------|-------------|
| `feat` | Nova feature ou endpoint |
| `fix` | Correção de bug |
| `refactor` | Reorganização sem mudar comportamento |
| `test` | Adicionar ou corrigir testes |
| `docs` | Documentação, READMEs, diagramas |
| `chore` | Dependências, config, Makefile |
| `perf` | Melhorias de performance |
| `ci` | Alterações ao pipeline CI/CD |

**Regra:** o título descreve o QUÊ. O body (opcional) explica o PORQUÊ.

```
feat(contact): add ILIKE search with trigram index

PostgreSQL full-table scan on ILIKE was timing out with >10k contacts.
GIN trigram index reduces search time from ~800ms to ~5ms.
```

## Razões

1. **Pedagogia** — cada branch conta a história de um módulo; `git log` é o índice do curso
2. **Rastreabilidade** — conventional commits permitem gerar CHANGELOG automaticamente
3. **Code review** — PRs com template garantem que cada módulo está documentado antes de mergear
4. **Realismo** — é o workflow usado em equipas profissionais

## Consequências

- O `main` só avança via PR — nunca por push direto
- Cada PR usa o template em `.github/pull_request_template.md`
- Tags semânticas (`v0.1.0`, `v0.2.0`...) marcam o fim de cada módulo
