# Clean Code — Antes e Depois no GoRM

Exemplos reais do refactoring aplicado neste módulo.
Cada secção é um princípio com código real do projecto.

---

## 1. Números Mágicos → Constantes Nomeadas

**Ficheiros:** `internal/activitylog/repository_mongo.go`, `internal/shared/events/events.go`

```go
// ❌ ANTES — o que significa 90 * 24 * 3600?
SetExpireAfterSeconds(90 * 24 * 3600)
bus := events.New(500, log)
context.WithTimeout(ctx, 10*time.Second)

// ✅ DEPOIS — lê-se como prosa
const logRetentionSecs = logRetentionDays * 24 * 60 * 60
const DefaultBufferSize = 500

SetExpireAfterSeconds(logRetentionSecs)
bus := events.New(events.DefaultBufferSize, log)
context.WithTimeout(ctx, indexTimeout)
```

**Regra:** se um número aparecer mais de uma vez, ou precisar de um comentário para ser compreendido, dá-lhe um nome.

---

## 2. Switch → Map (Open/Closed Principle)

**Ficheiro:** `internal/shared/validate/validate.go`

```go
// ❌ ANTES — adicionar nova tag = modificar a função
func fieldMessage(e validator.FieldError) string {
    switch e.Tag() {
    case "required": return "campo obrigatório"
    case "email":    return "email inválido"
    // cada nova tag abre este switch
    }
}

// ✅ DEPOIS — adicionar nova tag = uma linha no map
var messages = map[string]string{
    "required": "campo obrigatório",
    "email":    "email inválido",
    "uuid4":    "UUID inválido",    // ← sem tocar no código existente
}
// switch só para os casos que precisam de formatação dinâmica
```

**Regra:** se um switch cresce com novos cases ao longo do tempo, considera um map. O código deve estar aberto à extensão, fechado à modificação.

---

## 3. Strings Mágicas → Tipos Fortes

**Ficheiro:** `internal/activitylog/model.go`

```go
// ❌ ANTES — string literal espalhada
func entityTypeFromEventType(...) string {
    return "contact"  // e se escrevermos "Contact"? erro em runtime
}
repo.FindByEntity("contact", id, 50)  // sem verificação do compilador

// ✅ DEPOIS — tipo enumerado
type EntityType string
const (
    EntityContact EntityType = "contact"
    EntityLead    EntityType = "lead"
)
repo.FindByEntity(EntityContact, id, 50)  // compilador verifica
```

**Regra:** se uma string tem um conjunto finito de valores válidos, define um tipo. O compilador passa a ser o teu linter.

---

## 4. Comentários que Explicam O QUÊ → Remover

**Ficheiro:** `internal/contact/service.go`

```go
// ❌ ANTES — o nome já diz tudo
// CreateContactDTO é o input validado para criar um contacto.
type CreateContactDTO struct { ... }

// UpdateContactDTO é o input para atualizar — todos os campos opcionais.
// (os campos *string já comunicam "opcional")
type UpdateContactDTO struct { ... }

// ✅ DEPOIS — sem comentário redundante
type CreateContactDTO struct { ... }
type UpdateContactDTO struct { ... }
```

```go
// ✅ MANTER — explica o PORQUÊ, não óbvio do código
// handleEvent corre na goroutine do bus — o handler HTTP já respondeu.
// Falha silenciosa intencional: logs são "best effort", não críticos.
func (s *Service) handleEvent(...) { ... }
```

**Regra:** se podes apagar o comentário sem perder informação (porque o nome já a diz), apaga. Se o comentário explica uma decisão de design, uma limitação, ou algo que surpreenderia o leitor — mantém.

---

## 5. Early Return vs If/Else Aninhado

**Ficheiro:** `cmd/api/main.go`

```go
// ❌ ANTES — happy path enterrado no else
var mongoDB *mongo.Database
mongoDB, err = database.NewMongo(...)
if err != nil {
    log.Warn(...)
} else {
    log.Info(...)           // happy path dentro de else
    actSvc := ...
    actSvc.RegisterHandlers(bus)
}

// ✅ DEPOIS — early return, happy path à esquerda
func connectMongo(log) *mongo.Database {
    db, err := database.NewMongo(...)
    if err != nil {
        log.Warn(...)
        return nil          // erro: sai imediatamente
    }
    log.Info(...)
    return db               // sucesso: linha final, sem else
}
```

**Regra:** trata os casos de erro/guarda primeiro e sai. O happy path fica no nível principal, sem indentação adicional. A função `connectMongo` também tem o bónus de ter um nome descritivo.

---

## 6. Funções Pequenas — Uma Responsabilidade

**Ficheiro:** `internal/contact/service.go`

```go
// ❌ ANTES — Update faz duas coisas
func (s *Service) Update(id, dto) (*Contact, error) {
    contact, _ := s.repo.FindByID(id)
    // responsabilidade 1: orquestrar o update
    if dto.Name != nil { contact.Name = *dto.Name }    // responsabilidade 2:
    if dto.Phone != nil { contact.Phone = *dto.Phone } // aplicar campos
    if dto.Company != nil { contact.Company = *dto.Company }
    if dto.Notes != nil { contact.Notes = *dto.Notes }
    updated, _ := s.repo.Update(contact)
    s.bus.Publish(...)
    return updated, nil
}

// ✅ DEPOIS — cada função com uma responsabilidade
func (s *Service) Update(id, dto) (*Contact, error) {
    contact, _ := s.repo.FindByID(id)
    applyUpdates(contact, dto)  // delegado — Update não sabe os detalhes
    updated, _ := s.repo.Update(contact)
    s.bus.Publish(...)
    return updated, nil
}

func applyUpdates(c *Contact, dto UpdateContactDTO) {
    // única responsabilidade: aplicar campos opcionais
    if dto.Name != nil { c.Name = *dto.Name }
    // ...
}
```

**Regra:** se uma função tem mais de um "e" quando a describes em palavras, divide-a. "Busca o contacto **e** aplica os campos **e** guarda **e** publica" → três funções.

---

## Resumo dos Princípios Aplicados

| Princípio | Ficheiro(s) alterado(s) |
|-----------|------------------------|
| Sem números mágicos | `events.go`, `repository_mongo.go`, `main.go` |
| Map > switch crescente | `validate.go` |
| Tipos fortes > strings | `activitylog/model.go` e relacionados |
| Comentários só para PORQUÊ | `contact/service.go`, `activitylog/service.go` |
| Early return > if/else | `main.go` → `connectMongo` |
| Uma responsabilidade por função | `contact/service.go`, `task/service.go` |
