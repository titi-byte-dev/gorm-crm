package activitylog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

// Service tem duas responsabilidades distintas, separadas por metodo:
//   - Escrita: subscreve eventos e persiste logs (RegisterHandlers)
//   - Leitura: query de logs (GetByEntity, GetByUser)
// O EventMapper foi extraido para mapper.go — responsabilidade unica de mapeamento.
type Service struct {
	repo   Repository
	logger *slog.Logger
	mapper EventMapper
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger, mapper: EventMapper{}}
}

// RegisterHandlers subscreve todos os eventos relevantes no bus.
// Esta função é chamada uma vez no startup — depois corre em background.
//
// Padrão Observer em acção:
//   Emitter (ContactService, DealService...) → Event Bus → Observer (este service)
//   Os emitters não sabem que existe um ActivityLog.
//   O ActivityLog não sabe como os eventos são gerados.
//   Desacoplamento total.
func (s *Service) RegisterHandlers(bus events.Subscriber) {
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
	log := s.mapper.ToLog(event)
	if err := s.repo.Save(log); err != nil {
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

