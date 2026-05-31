package activitylog

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Log representa uma entrada no histórico de actividade do CRM.
//
// Porquê usar any (interface{}) para Payload?
// Cada tipo de evento tem dados diferentes:
//   ContactCreated → { name, email, company }
//   DealWon        → { title, value, contact_id }
//   TaskOverdue    → { id, due_date }
//
// Numa base de dados relacional, precisaríamos de uma tabela por tipo
// ou de uma coluna JSON. No MongoDB, cada documento pode ter estrutura
// diferente — é a flexibilidade que faz o Document Store valer a pena.
type Log struct {
	// primitive.ObjectID é o tipo nativo do MongoDB para _id
	// Diferente do uuid.UUID que usamos no PostgreSQL
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Action     string             `bson:"action"        json:"action"`
	EntityType string             `bson:"entity_type"   json:"entity_type"`
	EntityID   string             `bson:"entity_id"     json:"entity_id"`
	UserID     string             `bson:"user_id"       json:"user_id"`
	Payload    any                `bson:"payload"       json:"payload"`
	CreatedAt  time.Time          `bson:"created_at"    json:"created_at"`
}

// Repository define o contrato de acesso para ActivityLog.
// A implementação usa MongoDB; em testes usamos um mock em memória.
type Repository interface {
	Save(log *Log) error
	FindByEntity(entityType, entityID string, limit int) ([]*Log, error)
	FindByUser(userID string, limit int) ([]*Log, error)
}

// Filters para queries de logs.
type Filters struct {
	EntityType string
	EntityID   string
	UserID     string
	Limit      int
}

func (f *Filters) SetDefaults() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
}
