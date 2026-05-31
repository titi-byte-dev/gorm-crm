# 🎯 CHALLENGE — Módulo 18: Cloud & CI/CD

---

### Nível 1 — Verificar Version Injection

Confirma que os valores são mesmo injectados no binário:

```bash
make build
# ✅ bin/gorm-crm — version=v1.x.x-dirty commit=abc1234

# Inicia a app e chama o health check
./bin/gorm-crm &
curl -s localhost:8080/health | jq .
# Deve mostrar "version": "v1.x.x-dirty", "commit": "abc1234"
```

Agora adiciona `BuildTime` à resposta do `/health`:

```go
// cmd/api/main.go — no handler /health
"build_time": version.BuildTime,
```

> **Pergunta:** Por que é que `BuildTime` não aparece no `/health` por defeito? O que é que isso diz sobre a decisão de expor metadados de build em APIs públicas?

---

### Nível 2 — Job de Segurança no CI

Adiciona um job `security` ao `ci.yml`:

```yaml
security:
  name: Security Scan
  runs-on: ubuntu-latest
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-go@v5
      with:
        go-version: "1.22"
    - name: govulncheck
      run: |
        go install golang.org/x/vuln/cmd/govulncheck@latest
        govulncheck ./...
```

`govulncheck` verifica se o código usa funções de dependências com CVEs conhecidos.

> **Pergunta:** `govulncheck` é diferente de `go mod audit` ou de ferramentas como Snyk. Qual é a diferença principal? Em que tipo de vulnerabilidade cada um é melhor?

---

### Nível 3 — Deploy Automático via SSH

Estende o CD para fazer deploy após push da imagem:

```yaml
# No cd.yml, após o job release
deploy:
  name: Deploy
  needs: release
  runs-on: ubuntu-latest
  environment: production    # requer aprovação manual no GitHub
  steps:
    - name: Deploy via SSH
      uses: appleboy/ssh-action@v1
      with:
        host: ${{ secrets.DEPLOY_HOST }}
        username: ${{ secrets.DEPLOY_USER }}
        key: ${{ secrets.DEPLOY_KEY }}
        script: |
          docker pull ghcr.io/${{ github.repository }}:${{ github.ref_name }}
          docker stop gorm-crm || true
          docker run -d --name gorm-crm \
            --env-file /etc/gorm-crm/prod.env \
            -p 8080:8080 \
            ghcr.io/${{ github.repository }}:${{ github.ref_name }}
```

> **Atenção:** O `environment: production` com `required reviewers` no GitHub garante que o deploy não acontece sem aprovação — mesmo que o CI passe.

---

## Perguntas de reflexão

1. **CI vs CD:** O CI corre em cada push. O CD só em tags. Por que esta separação? O que aconteceria se o CD corresse em cada push para main?

2. **Imagem mínima:** O Dockerfile usa `alpine:3.19` como runtime. `scratch` seria ainda menor (~0MB overhead). Quando usarias `scratch`? Qual é o trade-off com `alpine`?

3. **Secrets:** O `GITHUB_TOKEN` é automático. `DEPLOY_KEY` e `DEPLOY_HOST` precisam de ser configurados manualmente. Qual é a diferença de risco entre os dois? Como rotacionas `DEPLOY_KEY` sem downtime?

4. **Version drift:** Se o CD falhar a meio, a imagem está publicada mas o deploy não aconteceu. Como detectas este estado? Como resolves?

---

> 🏁 **Fim do curso GoRM.** Construíste um CRM completo com Go, GORM, PostgreSQL, MongoDB, testes automatizados, design patterns, refactoring, performance e pipeline CI/CD.
>
> O próximo passo é o frontend — ver [FRONTEND.md](FRONTEND.md) para o plano.
