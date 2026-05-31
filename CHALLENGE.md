# 🎯 CHALLENGE — Módulo 04: Git Workflow

---

## Desafios

### Nível 1 — Obrigatório

**Ativa o commit template e faz um commit "a sério"**

```bash
make setup
# Edita qualquer ficheiro (ex: adiciona um comentário ao README)
git add -p
git commit     # abre o editor com o template
```

Escreve um commit com tipo, escopo e descrição válidos. Experimenta preencher o body também.

---

### Nível 2 — Exploração

**Dois commits a partir de um ficheiro com duas mudanças**

1. Abre `internal/contact/service.go`
2. Faz duas mudanças não relacionadas (ex: adiciona um método + muda uma mensagem de erro)
3. Usa `git add -p` para separar as duas mudanças
4. Faz dois commits distintos, cada um com o seu tipo e escopo

Verifica com `git log --oneline` que tens dois commits limpos.

---

### Nível 3 — Investigação

**Cria uma branch de desafio e abre um PR**

```bash
git checkout branch-03-sql
git checkout -b meu-desafio-m03-stats
```

Implementa o endpoint `GET /api/v1/contacts/stats` do Challenge do M03.

Depois:
```bash
git push origin meu-desafio-m03-stats
gh pr create --base branch-03-sql --head meu-desafio-m03-stats \
  --title "challenge(m03): add contacts stats endpoint"
```

---

## Perguntas de reflexão

1. Qual a diferença entre `git merge` e `git rebase`? Quando usarias cada um?
2. O que é um "squash merge" e que vantagens tem para manter o histórico limpo?
3. Porque é que `git add .` pode ser perigoso em projetos reais?

---

> Módulo seguinte: [branch-05-rest-api](https://github.com/titi-byte-dev/gorm-crm/tree/branch-05-rest-api) — API REST completa, middlewares e paginação
