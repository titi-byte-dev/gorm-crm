package lead

import (
	"fmt"

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

type CreateLeadInput struct {
	Title     string    `json:"title"      validate:"required,min=2,max=200"`
	Value     float64   `json:"value"      validate:"min=0"`
	ContactID uuid.UUID `json:"contact_id" validate:"required"`
}

type UpdateLeadInput struct {
	Title  *string  `json:"title"  validate:"omitempty,min=2,max=200"`
	Value  *float64 `json:"value"  validate:"omitempty,min=0"`
	Status *Status  `json:"status" validate:"omitempty,oneof=new contacted qualified lost"`
}

func (s *Service) Create(ownerID uuid.UUID, dto CreateLeadInput) (*Lead, error) {
	value, err := valueobject.ParseMoney(dto.Value)
	if err != nil {
		return nil, fmt.Errorf("create lead: %w", err)
	}
	lead := &Lead{
		Title:     dto.Title,
		Value:     value,
		Status:    StatusNew,
		ContactID: dto.ContactID,
		OwnerID:   ownerID,
	}
	saved, err := s.repo.Save(lead)
	if err != nil {
		return nil, fmt.Errorf("create lead: %w", err)
	}
	s.bus.Publish(events.Event{Type: events.LeadCreated, Payload: saved, UserID: ownerID.String()})
	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID) (*Lead, error) {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get lead: %w", err)
	}
	return lead, nil
}

func (s *Service) List(ownerID uuid.UUID, filters Filters) ([]*Lead, int64, error) {
	return s.repo.FindAll(ownerID, filters)
}

func (s *Service) UpdateStatus(id uuid.UUID, newStatus Status) (*Lead, error) {
	lead, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update lead status: %w", err)
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
		s.bus.Publish(events.Event{Type: events.LeadLost, Payload: updated, UserID: lead.OwnerID.String()})
	} else if newStatus == StatusQualified {
		s.bus.Publish(events.Event{Type: events.LeadConverted, Payload: updated, UserID: lead.OwnerID.String()})
	}

	return updated, nil
}

func (s *Service) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}
