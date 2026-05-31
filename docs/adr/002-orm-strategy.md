# ADR-002 — Estratégia de ORM vs SQL Puro

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Módulo 03 (SQL) e Módulo 17 (Performance)

---

## Contexto

Em Go existem três abordagens principais para acesso a dados:

| Abordagem | Biblioteca | Prós | Contras |
|-----------|-----------|------|---------|
| ORM completo | GORM | Produtividade alta, migrations, associations | Magic escondida, queries menos controladas |
| SQL semi-nativo | sqlx | Controlo total do SQL, mapeamento simples | Mais verboso, migrations manuais |
| SQL puro | `database/sql` | Zero abstração, máximo controlo | Muito verboso, scan manual |

## Decisão

**Abordagem progressiva:**

1. **Módulos 03–16:** Usar **GORM** — o estudante foca-se nos conceitos de backend, não no SQL boilerplate
2. **Módulo 17 (Performance):** Introduzir **sqlx** para demonstrar queries otimizadas e o custo das abstrações

## Razões

1. **Didático** — o estudante aprende primeiro a produzir, depois a otimizar
2. **Real world** — a maioria dos projetos começa com ORM e otimiza partes críticas com SQL mais direto
3. **Trade-offs visíveis** — ao introduzir sqlx no módulo de performance, o estudante vê concretamente o que GORM abstrai e quando essa abstração tem custo

## Consequências

- GORM gera queries menos eficientes em alguns casos — isso é um ponto de ensino intencional no módulo 17
- As interfaces de repositório (Repository Pattern) tornam a troca transparente — outro ponto de ensino
