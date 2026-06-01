package deal

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
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

func (s *Service) Create(rctx ctxutil.RequestCtx, dto CreateDealDTO) (*Deal, error) {
	deal := &Deal{
		Title:     dto.Title,
		Value:     dto.Value,
		Stage:     StageProposal,
		ContactID: dto.ContactID,
		LeadID:    dto.LeadID,
		OwnerID:   rctx.UserID,
		TenantID:  rctx.TenantID,
	}
	saved, err := s.repo.Save(deal)
	if err != nil {
		return nil, fmt.Errorf("create deal: %w", err)
	}
	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID, rctx ctxutil.RequestCtx) (*Deal, error) {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := s.checkAccess(deal, rctx); err != nil {
		return nil, err
	}
	return deal, nil
}

func (s *Service) List(rctx ctxutil.RequestCtx, filters Filters) ([]*Deal, int64, error) {
	return s.repo.FindAll(rctx.TenantID, rctx.UserID, rctx.IsManager(), filters)
}

func (s *Service) MoveStage(id uuid.UUID, rctx ctxutil.RequestCtx, newStage Stage) (*Deal, error) {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("move stage: %w", err)
	}
	if err := s.checkAccess(deal, rctx); err != nil {
		return nil, err
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
		s.bus.Publish(events.Event{Type: evtType, Payload: updated, UserID: rctx.UserID.String()})
	}

	return updated, nil
}

func (s *Service) Delete(id uuid.UUID, rctx ctxutil.RequestCtx) error {
	deal, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete deal: %w", err)
	}
	if err := s.checkAccess(deal, rctx); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) checkAccess(d *Deal, rctx ctxutil.RequestCtx) error {
	if d.TenantID != rctx.TenantID {
		return fmt.Errorf("deal %s: %w", d.ID, sharederrors.ErrNotFound)
	}
	if !rctx.IsManager() && d.OwnerID != rctx.UserID {
		return fmt.Errorf("deal %s: %w", d.ID, sharederrors.ErrNotFound)
	}
	return nil
}
