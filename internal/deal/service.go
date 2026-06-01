package deal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

type Service struct {
	repo Repository
	bus  *events.Bus
}

func NewService(repo Repository, bus *events.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

type CreateDealDTO struct {
	Title     string     `json:"title"      validate:"required,min=2,max=200"`
	Value     float64    `json:"value"      validate:"min=0"`
	ContactID uuid.UUID  `json:"contact_id" validate:"required"`
	LeadID    *uuid.UUID `json:"lead_id"`
}

func (s *Service) Create(ownerID uuid.UUID, dto CreateDealDTO) (*Deal, error) {
	deal := &Deal{
		Title:     dto.Title,
		Value:     dto.Value,
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

func (s *Service) GetByID(id, ownerID uuid.UUID) (*Deal, error) {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if deal.OwnerID != ownerID {
		return nil, fmt.Errorf("deal %s: %w", id, sharederrors.ErrNotFound)
	}
	return deal, nil
}

func (s *Service) List(ownerID uuid.UUID, filters Filters) ([]*Deal, int64, error) {
	return s.repo.FindAll(ownerID, filters)
}

func (s *Service) MoveStage(id, ownerID uuid.UUID, newStage Stage) (*Deal, error) {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("move stage: %w", err)
	}
	if deal.OwnerID != ownerID {
		return nil, fmt.Errorf("deal %s: %w", id, sharederrors.ErrNotFound)
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

	evtType := events.DealLost
	if newStage == StageWon {
		evtType = events.DealWon
	}
	if newStage.IsClosed() {
		s.bus.Publish(events.Event{Type: evtType, Payload: updated, UserID: deal.OwnerID.String()})
	}

	return updated, nil
}

func (s *Service) Delete(id, ownerID uuid.UUID) error {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}
	if deal.OwnerID != ownerID {
		return fmt.Errorf("deal %s: %w", id, sharederrors.ErrNotFound)
	}
	return s.repo.Delete(id)
}
