package contact

import (
	"time"

	"github.com/google/uuid"
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
	TenantID  uuid.UUID `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository define o contrato de acesso a dados para Contact.
// A implementação concreta (PostgreSQL, mock para testes) fica noutro ficheiro.
type Repository interface {
	FindByID(id uuid.UUID) (*Contact, error)
	FindAll(tenantID, ownerID uuid.UUID, isManager bool, filters Filters) ([]*Contact, int64, error)
	FindByEmail(email string) (*Contact, error)
	Save(contact *Contact) (*Contact, error)
	Update(contact *Contact) (*Contact, error)
	Delete(id uuid.UUID) error
}

// Filters encapsula os parâmetros de pesquisa e paginação.
// Usar uma struct em vez de parâmetros avulsos torna a assinatura estável
// — adicionar um novo filtro não quebra os callers existentes.
type Filters struct {
	Search  string
	Company string
	Page    int
	Limit   int
	SortBy  string
	SortDir string
}

func (f *Filters) SetDefaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.SortBy == "" {
		f.SortBy = "created_at"
	}
	if f.SortDir == "" {
		f.SortDir = "desc"
	}
}

func (f *Filters) Offset() int {
	return (f.Page - 1) * f.Limit
}
