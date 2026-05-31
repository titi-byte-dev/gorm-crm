package activitylog

import "github.com/titi-byte-dev/gorm-crm/internal/shared/events"

// EventMapper tem uma unica responsabilidade: converter events.Event em Log.
// SRP — Single Responsibility Principle: uma razao para mudar.
// Se a logica de mapeamento mudar, so este ficheiro e alterado.
// Se a logica de persistencia mudar, so service.go e alterado.
type EventMapper struct{}

func (m EventMapper) ToLog(event events.Event) *Log {
	entityType, entityID := m.extractEntity(event)
	return &Log{
		Action:     string(event.Type),
		UserID:     event.UserID,
		Payload:    event.Payload,
		EntityType: entityType,
		EntityID:   entityID,
	}
}

func (m EventMapper) extractEntity(event events.Event) (EntityType, string) {
	entityType := m.entityTypeFromEventType(event.Type)
	var entityID string
	switch p := event.Payload.(type) {
	case interface{ GetID() string }:
		entityID = p.GetID()
	case map[string]string:
		if id, ok := p["id"]; ok {
			entityID = id
		}
	}
	return entityType, entityID
}

func (m EventMapper) entityTypeFromEventType(et events.EventType) EntityType {
	switch et {
	case events.ContactCreated, events.ContactUpdated, events.ContactDeleted:
		return EntityContact
	case events.LeadCreated, events.LeadConverted, events.LeadLost:
		return EntityLead
	case events.DealWon, events.DealLost:
		return EntityDeal
	case events.TaskOverdue:
		return EntityTask
	default:
		return EntityUnknown
	}
}
