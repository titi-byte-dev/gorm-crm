package contact

import (
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
)

// Contact representa um contacto no CRM.
type Contact struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	Company   string    `json:"company,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository define o contrato de acesso a dados para Contact.
// A implementação concreta (PostgreSQL, mock para testes) fica noutro ficheiro.
type Repository interface {
	FindByID(id uuid.UUID) (*Contact, error)
	FindAll(ownerID uuid.UUID, filters Filters) ([]*Contact, int64, error)
	FindByEmail(email string) (*Contact, error)
	Save(contact *Contact) (*Contact, error)
	Update(contact *Contact) (*Contact, error)
	Delete(id uuid.UUID) error
}

// Filters encapsula os parâmetros de pesquisa e paginação para contactos.
// Embebe pagination.Base para herdar Page, Limit, SortBy, SortDir e os métodos
// Normalize/Offset — composição em vez de duplicação.
type Filters struct {
	pagination.Base
	Search  string
	Company string
}

func (f *Filters) SetDefaults() {
	f.Normalize("created_at")
}
