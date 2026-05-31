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

// Reader define operacoes de leitura sobre Contact.
// Um servico de relatorios pode receber apenas Reader — superficie minima.
type Reader interface {
	FindByID(id uuid.UUID) (*Contact, error)
	FindAll(ownerID uuid.UUID, filters Filters) ([]*Contact, int64, error)
	FindByEmail(email string) (*Contact, error)
}

// Writer define operacoes de escrita sobre Contact.
type Writer interface {
	Save(contact *Contact) (*Contact, error)
	Update(contact *Contact) (*Contact, error)
	Delete(id uuid.UUID) error
}

// Repository e a composicao de Reader e Writer — contrato completo para o Service.
// Interface embedding: Repository inclui todos os metodos de Reader e Writer.
type Repository interface {
	Reader
	Writer
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
