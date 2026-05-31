package contact

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

type Service struct {
	repo  Repository
	bus   events.Publisher
	rules []Rule
}

// NewService aceita regras de negocio opcionais via variadic.
// OCP: para adicionar uma nova regra, passa-a aqui — Service.Create nao muda.
// Sem regras passadas, aplica UniqueEmailRule por defeito.
func NewService(repo Repository, bus events.Publisher, rules ...Rule) *Service {
	if len(rules) == 0 {
		rules = []Rule{UniqueEmailRule{}}
	}
	return &Service{repo: repo, bus: bus, rules: rules}
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

func (s *Service) Create(ownerID uuid.UUID, dto CreateContactDTO) (*Contact, error) {
	for _, rule := range s.rules {
		if err := rule.Validate(s.repo, dto); err != nil {
			return nil, fmt.Errorf("contact rule: %w", err)
		}
	}

	contact := &Contact{
		Name:    dto.Name,
		Email:   dto.Email,
		Phone:   dto.Phone,
		Company: dto.Company,
		Notes:   dto.Notes,
		OwnerID: ownerID,
	}

	saved, err := s.repo.Save(contact)
	if err != nil {
		return nil, fmt.Errorf("create contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactCreated,
		Payload: saved,
		UserID:  ownerID.String(),
	})

	return saved, nil
}

func (s *Service) GetByID(id uuid.UUID) (*Contact, error) {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("get contact: %w", err)
	}
	return contact, nil
}

func (s *Service) List(ownerID uuid.UUID, filters Filters) (Contacts, int64, error) {
	contacts, total, err := s.repo.FindAll(ownerID, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("list contacts: %w", err)
	}
	return contacts, total, nil
}

func (s *Service) Update(id uuid.UUID, dto UpdateContactDTO) (*Contact, error) {
	contact, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("update contact: %w", err)
	}

	applyUpdates(contact, dto)

	updated, err := s.repo.Update(contact)
	if err != nil {
		return nil, fmt.Errorf("update contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactUpdated,
		Payload: updated,
		UserID:  contact.OwnerID.String(),
	})

	return updated, nil
}

// applyUpdates aplica os campos opcionais do DTO ao contacto.
// Ponteiro nil significa "não alterar" — só actualiza campos enviados.
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

func (s *Service) Delete(id uuid.UUID) error {
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("delete contact: %w", err)
	}

	s.bus.Publish(events.Event{
		Type:    events.ContactDeleted,
		Payload: map[string]string{"id": id.String()},
	})

	return nil
}
