package activitylog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

// RegisterHandlers subscreve todos os eventos relevantes no bus.
// Esta função é chamada uma vez no startup — depois corre em background.
//
// Padrão Observer em acção:
//   Emitter (ContactService, DealService...) → Event Bus → Observer (este service)
//   Os emitters não sabem que existe um ActivityLog.
//   O ActivityLog não sabe como os eventos são gerados.
//   Desacoplamento total.
func (s *Service) RegisterHandlers(bus *events.Bus) {
	eventTypes := []events.EventType{
		events.ContactCreated,
		events.ContactUpdated,
		events.ContactDeleted,
		events.LeadCreated,
		events.LeadConverted,
		events.LeadLost,
		events.DealWon,
		events.DealLost,
		events.TaskOverdue,
	}
	for _, et := range eventTypes {
		bus.Subscribe(et, s.handleEvent)
	}
}

// handleEvent corre na goroutine do bus — o handler HTTP já respondeu ao utilizador.
// Falha silenciosa intencional: logs são "best effort", não críticos para o negócio.
func (s *Service) handleEvent(ctx context.Context, event events.Event) {
	log := &Log{
		Action:  string(event.Type),
		UserID:  event.UserID,
		Payload: event.Payload,
	}

	// Extrai entity_type e entity_id do tipo de evento
	entityType, entityID := entityFromEvent(event)
	log.EntityType = entityType
	log.EntityID = entityID

	if err := s.repo.Save(log); err != nil {
		// Log do erro mas não propaga — falhar um log não deve quebrar a operação principal
		s.logger.Error("failed to save activity log", "event", event.Type, "error", err)
	}
}

func (s *Service) GetByEntity(entityType EntityType, entityID string, limit int) ([]*Log, error) {
	logs, err := s.repo.FindByEntity(entityType, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("get activity logs: %w", err)
	}
	return logs, nil
}

func (s *Service) GetByUser(userID string, limit int) ([]*Log, error) {
	logs, err := s.repo.FindByUser(userID, limit)
	if err != nil {
		return nil, fmt.Errorf("get user activity: %w", err)
	}
	return logs, nil
}

// entityFromEvent extrai o tipo e ID da entidade afectada pelo evento.
func entityFromEvent(event events.Event) (entityType EntityType, entityID string) {
	entityType = entityTypeFromEventType(event.Type)
	// O payload é any — fazemos type assertion para extrair o ID
	// Esta é a "taxa" pelo uso de interface{} — precisamos de lidar com cada tipo
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

func entityTypeFromEventType(et events.EventType) EntityType {
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
