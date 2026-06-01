package lead

import (
	"fmt"

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

type CreateLeadDTO struct {
	Title     string    `json:"title"      validate:"required,min=2,max=200"`
	Value     float64   `json:"value"      validate:"min=0"`
	ContactID uuid.UUID `json:"contact_id" validate:"required"`
}

type UpdateLeadDTO struct {
	Title  *string  `json:"title"  validate:"omitempty,min=2,max=200"`
	Value  *float64 `json:"value"  validate:"omitempty,min=0"`
	Status *Status  `json:"status" validate:"omitempty,oneof=new contacted qualified lost"`
}

func (s *Service) Create(rctx ctxutil.RequestCtx, dto CreateLeadDTO) (*Lead, error) {
	lead := &Lead{
		Title:     dto.Title,
		Value:     dto.Value,
		Status:    StatusNew,
		ContactID: dto.ContactID,
		OwnerID:   rctx.UserID,
		TenantID:  rctx.TenantID,
	}
	saved, err := s.repo.Save(lead)
	if err != nil {
		return nil, fmt.Errorf("create lead: %w", err)
	}
	s.bus.Publish(events.Event{Type: events.LeadCreated, Payload: saved, UserID: rctx.UserID.String()})
	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID, rctx ctxutil.RequestCtx) (*Lead, error) {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get lead: %w", err)
	}
	if err := s.checkAccess(lead, rctx); err != nil {
		return nil, err
	}
	return lead, nil
}

func (s *Service) List(rctx ctxutil.RequestCtx, filters Filters) ([]*Lead, int64, error) {
	return s.repo.FindAll(rctx.TenantID, rctx.UserID, rctx.IsManager(), filters)
}

func (s *Service) UpdateStatus(id uuid.UUID, rctx ctxutil.RequestCtx, newStatus Status) (*Lead, error) {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update lead status: %w", err)
	}
	if err := s.checkAccess(lead, rctx); err != nil {
		return nil, err
	}

	if !lead.Status.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf("cannot transition from %s to %s: %w",
			lead.Status, newStatus, sharederrors.ErrValidation)
	}

	lead.Status = newStatus
	updated, err := s.repo.Update(lead)
	if err != nil {
		return nil, fmt.Errorf("update lead: %w", err)
	}

	if newStatus == StatusLost {
		s.bus.Publish(events.Event{Type: events.LeadLost, Payload: updated, UserID: rctx.UserID.String()})
	} else if newStatus == StatusQualified {
		s.bus.Publish(events.Event{Type: events.LeadConverted, Payload: updated, UserID: rctx.UserID.String()})
	}

	return updated, nil
}

func (s *Service) Delete(id uuid.UUID, rctx ctxutil.RequestCtx) error {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete lead: %w", err)
	}
	if err := s.checkAccess(lead, rctx); err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) checkAccess(l *Lead, rctx ctxutil.RequestCtx) error {
	if l.TenantID != rctx.TenantID {
		return fmt.Errorf("lead %s: %w", l.ID, sharederrors.ErrNotFound)
	}
	if !rctx.IsManager() && l.OwnerID != rctx.UserID {
		return fmt.Errorf("lead %s: %w", l.ID, sharederrors.ErrNotFound)
	}
	return nil
}
