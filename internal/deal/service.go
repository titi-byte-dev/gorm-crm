package deal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	"github.com/titi-byte-dev/gorm-crm/pkg/valueobject"
)

type Service struct {
	repo Repository
	bus  events.Publisher
}

func NewService(repo Repository, bus events.Publisher) *Service {
	return &Service{repo: repo, bus: bus}
}

type CreateDealInput struct {
	Title     string     `json:"title"      validate:"required,min=2,max=200"`
	Value     float64    `json:"value"      validate:"min=0"`
	ContactID uuid.UUID  `json:"contact_id" validate:"required"`
	LeadID    *uuid.UUID `json:"lead_id"`
}

func (s *Service) Create(ownerID uuid.UUID, dto CreateDealInput) (*Deal, error) {
	value, err := valueobject.ParseMoney(dto.Value)
	if err != nil {
		return nil, fmt.Errorf("create deal: %w", err)
	}
	deal := &Deal{
		Title:     dto.Title,
		Value:     value,
		Stage:     StageProposal,
		ContactID: dto.ContactID,
		LeadID:    dto.LeadID,
		OwnerID:   ownerID,
	}
	saved, err := s.repo.Save(deal)
	if err != nil {
		return nil, fmt.Errorf("create deal: %w", err)
	}
	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID) (*Deal, error) {
	return s.repo.FindByID(id)
}

func (s *Service) List(ownerID uuid.UUID, filters Filters) ([]*Deal, int64, error) {
	return s.repo.FindAll(ownerID, filters)
}

func (s *Service) MoveStage(id uuid.UUID, newStage Stage) (*Deal, error) {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("move stage: %w", err)
	}
	if !deal.Stage.CanTransitionTo(newStage) {
		return nil, fmt.Errorf("cannot move from %s to %s: %w",
			deal.Stage, newStage, sharederrors.ErrValidation)
	}

	deal.Stage = newStage
	if newStage.IsClosed() {
		now := time.Now()
		deal.ClosedAt = &now
	}

	updated, err := s.repo.Update(deal)
	if err != nil {
		return nil, fmt.Errorf("update deal: %w", err)
	}

	if newStage.IsClosed() {
		s.bus.Publish(events.Event{Type: closedEventType(newStage), Payload: updated, UserID: deal.OwnerID.String()})
	}

	return updated, nil
}

func closedEventType(stage Stage) events.EventType {
	if stage == StageWon {
		return events.DealWon
	}
	return events.DealLost
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
