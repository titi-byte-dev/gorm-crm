package contact

import (
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
)

// Contacts e uma coleccao de primeiro nivel — nao e apenas um slice.
// Object Calisthenics — Regra 4: First-class collections.
// Encapsula a iteracao e expoe comportamento de dominio em vez de expor o slice.
type Contacts []*Contact

// IDs devolve os UUIDs de todos os contactos.
func (cs Contacts) IDs() []uuid.UUID {
	ids := make([]uuid.UUID, len(cs))
	for i, c := range cs {
		ids[i] = c.ID
	}
	return ids
}

// FilterByCompany devolve apenas os contactos da empresa dada.
func (cs Contacts) FilterByCompany(company string) Contacts {
	var result Contacts
	for _, c := range cs {
		if c.Company == company {
			result = append(result, c)
		}
	}
	return result
}

// HasEmail verifica se algum contacto tem o email dado.
func (cs Contacts) HasEmail(email string) bool {
	for _, c := range cs {
		if c.Email == email {
			return true
		}
	}
	return false
}

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
	FindAll(ownerID uuid.UUID, filters Filters) (Contacts, int64, error)
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
