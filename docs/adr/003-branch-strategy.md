# ADR-003 — Estratégia de Branches do Curso

**Data:** 2026-05-31
**Estado:** ✅ Aceite
**Contexto:** Estrutura pedagógica do repositório

---

## Contexto

O Git é a estrutura de navegação do curso. Há duas abordagens possíveis:

| Abordagem | Descrição | Prós | Contras |
|-----------|-----------|------|---------|
| **Cumulativa** | Cada branch parte da anterior e acrescenta | Código cresce naturalmente, contexto sempre presente | Diff por módulo requer `git diff` |
| Isolada | Cada branch começa do zero para o módulo | Cada módulo é independente | Repetição de boilerplate, sem continuidade |

## Decisão

**Branches cumulativas** — cada branch parte do estado final da branch anterior.

```
main ← branch-01 ← branch-02 ← branch-03 ← ... ← branch-18
```

O `main` reflete sempre o estado mais completo/avançado (merge de cada módulo quando finalizado).

## Razões

1. **A app conta uma história** — o histórico Git é a narrativa do curso
2. **Contexto sempre presente** — ao fazer `git checkout branch-07-mvc`, tens todo o código dos módulos anteriores disponível
3. **Realismo** — é exatamente como se trabalha em projetos reais
4. **Progressão visível** — `git log --oneline` mostra a jornada completa

## Navegação

```bash
# Ver o que mudou num módulo específico
git diff branch-04-git-workflow..branch-05-rest-api

# Ver apenas os ficheiros alterados num módulo
git diff --name-only branch-06-auth..branch-07-mvc-layers

# Ver o histórico do módulo atual
git log --oneline branch-07-mvc-layers
```

## Convenção de tags

Cada módulo finalizado recebe uma tag semântica:

```
v0.1.0  → branch-01-setup
v0.2.0  → branch-02-go-fundamentos
...
v0.8.0  → branch-08-docker  (🏆 Júnior)
...
v1.5.0  → branch-15-patterns (🎯 Pleno)
...
v2.0.0  → branch-18-cloud-cicd (🎓 Sénior)
```
