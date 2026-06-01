package contact

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

type CreateContactDTO struct {
	Name    string `json:"name"    validate:"required,min=2,max=100"`
	Email   string `json:"email"   validate:"required,email"`
	Phone   string `json:"phone"   validate:"omitempty,max=20"`
	Company string `json:"company" validate:"omitempty,max=100"`
	Notes   string `json:"notes"   validate:"omitempty,max=1000"`
}

type UpdateContactDTO struct {
	Name    *string `json:"name"    validate:"omitempty,min=2,max=100"`
	Phone   *string `json:"phone"   validate:"omitempty,max=20"`
	Company *string `json:"company" validate:"omitempty,max=100"`
	Notes   *string `json:"notes"   validate:"omitempty,max=1000"`
}

func (s *Service) Create(rctx ctxutil.RequestCtx, dto CreateContactDTO) (*Contact, error) {
	// Regra de negócio: email único por owner (não global)
	existing, err := s.repo.FindByEmail(dto.Email)
	if err == nil && existing != nil && existing.OwnerID == rctx.UserID {
		return nil, fmt.Errorf("email already exists: %w", sharederrors.ErrConflict)
	}

	contact := &Contact{
		Name:     dto.Name,
		Email:    dto.Email,
		Phone:    dto.Phone,
		Company:  dto.Company,
		Notes:    dto.Notes,
		OwnerID:  rctx.UserID,
		TenantID: rctx.TenantID,
	}

	saved, err := s.repo.Save(contact)
	if err != nil {
		return nil, fmt.Errorf("create contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactCreated,
		Payload: saved,
		UserID:  rctx.UserID.String(),
	})

	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID, rctx ctxutil.RequestCtx) (*Contact, error) {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get contact: %w", err)
	}
	if err := s.checkAccess(contact, rctx); err != nil {
		return nil, err
	}
	return contact, nil
}

func (s *Service) List(rctx ctxutil.RequestCtx, filters Filters) ([]*Contact, int64, error) {
	contacts, total, err := s.repo.FindAll(rctx.TenantID, rctx.UserID, rctx.IsManager(), filters)
	if err != nil {
		return nil, 0, fmt.Errorf("list contacts: %w", err)
	}
	return contacts, total, nil
}

func (s *Service) Update(id uuid.UUID, rctx ctxutil.RequestCtx, dto UpdateContactDTO) (*Contact, error) {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update contact: %w", err)
	}
	if err := s.checkAccess(contact, rctx); err != nil {
		return nil, err
	}

	applyUpdates(contact, dto)

	updated, err := s.repo.Update(contact)
	if err != nil {
		return nil, fmt.Errorf("update contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactUpdated,
		Payload: updated,
		UserID:  rctx.UserID.String(),
	})

	return updated, nil
}

func (s *Service) Delete(id uuid.UUID, rctx ctxutil.RequestCtx) error {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}
	if err := s.checkAccess(contact, rctx); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactDeleted,
		Payload: map[string]string{"id": id.String()},
		UserID:  rctx.UserID.String(),
	})

	return nil
}

// checkAccess verifica que o utilizador tem permissão para aceder ao contacto.
// Manager/admin: basta o tenant_id coincidir.
// Seller: exige também que seja o owner.
func (s *Service) checkAccess(c *Contact, rctx ctxutil.RequestCtx) error {
	if c.TenantID != rctx.TenantID {
		return fmt.Errorf("contact %s: %w", c.ID, sharederrors.ErrNotFound)
	}
	if !rctx.IsManager() && c.OwnerID != rctx.UserID {
		return fmt.Errorf("contact %s: %w", c.ID, sharederrors.ErrNotFound)
	}
	return nil
}

func applyUpdates(c *Contact, dto UpdateContactDTO) {
	if dto.Name != nil {
		c.Name = *dto.Name
	}
	if dto.Phone != nil {
		c.Phone = *dto.Phone
	}
	if dto.Company != nil {
		c.Company = *dto.Company
	}
	if dto.Notes != nil {
		c.Notes = *dto.Notes
	}
}
