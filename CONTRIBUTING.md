# Guia de Contribuição & Navegação

---

## Para estudantes

### Começar do zero

```bash
git clone https://github.com/titi-byte-dev/gorm-crm.git
cd gorm-crm
git checkout branch-01-setup
cat README.md        # lê o contexto do módulo
make run             # corre a app
```

### Navegar entre módulos

```bash
# Ver todos os módulos
git branch -a | grep branch-

# Ir para qualquer módulo
git checkout branch-07-mvc-layers

# Ver o que mudou neste módulo em relação ao anterior
git diff branch-06-auth..branch-07-mvc-layers
```

### Trabalhar no desafio de um módulo

```bash
# Estás no branch-05-rest-api, queres fazer o desafio
git checkout -b meu-desafio-m05
# ... faz o teu trabalho ...
git add -p           # adiciona com revisão
git commit -m "feat: adicionar endpoint de pesquisa de contactos"
```

---

## Conventional Commits

Este repositório usa [Conventional Commits](https://www.conventionalcommits.org/):

```
<tipo>(<escopo opcional>): <descrição curta>

Exemplos:
feat(contact): add search endpoint with filters
fix(auth): handle expired refresh token correctly
refactor(deal): extract stage transition logic to service
test(contact): add integration tests for repository
docs(m07): add MVC diagram and challenge
chore: update go dependencies
```

| Tipo | Quando usar |
|------|-------------|
| `feat` | Nova feature ou endpoint |
| `fix` | Correção de bug |
| `refactor` | Reorganização sem mudar comportamento |
| `test` | Adicionar ou corrigir testes |
| `docs` | Documentação, diagramas, READMEs |
| `chore` | Dependências, configurações, Makefile |
| `perf` | Melhorias de performance |

---

## Pull Requests

Cada módulo finalizado é mergeado em `main` via PR. O título do PR segue o formato:

```
✅ Módulo 07 — Arquitetura MVC em Camadas
```

O corpo do PR inclui:
- O que foi construído
- Conceitos abordados
- Diagrama de contexto (se aplicável)
- Checklist de conclusão

---

## Para futuros instrutores

### Estrutura pedagógica

Cada branch é auto-suficiente para ensino:

1. `README.md` — contexto completo do módulo
2. `CHALLENGE.md` — exercício para os alunos
3. `docs/` — diagramas e decisões de design
4. Código Go — limpo, comentado apenas onde necessário

### Adaptar para uma turma

```bash
# Criar fork do repositório
# Criar branches de exercício sem a solução
git checkout branch-05-rest-api
git checkout -b exercise-05-rest-api
# Remover a solução, deixar os testes e o esqueleto
git push origin exercise-05-rest-api
```

### Princípios pedagógicos aplicados

- **Autoconstrução** — o aluno constrói, não copia
- **80/20** — foco no que aparece em 80% dos projetos reais
- **Contexto sempre presente** — cada branch tem o código anterior como base
- **Desafio prático** — cada módulo tem um `CHALLENGE.md`
- **Diagramas primeiro** — o aluno entende a estrutura antes de ver o código
