package lead

import (
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
)

// Status representa o estado de um lead no pipeline.
type Status string

const (
	StatusNew        Status = "new"
	StatusContacted  Status = "contacted"
	StatusQualified  Status = "qualified"
	StatusLost       Status = "lost"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusNew, StatusContacted, StatusQualified, StatusLost:
		return true
	}
	return false
}

// CanTransitionTo define as transições de estado válidas.
// Isto é o início do padrão State — expandido no Módulo 15.
func (s Status) CanTransitionTo(next Status) bool {
	transitions := map[Status][]Status{
		StatusNew:       {StatusContacted, StatusLost},
		StatusContacted: {StatusQualified, StatusLost},
		StatusQualified: {StatusLost},
		StatusLost:      {},
	}
	for _, allowed := range transitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

// Lead representa um potencial negócio no CRM.
type Lead struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Value     float64   `json:"value"`
	Status    Status    `json:"status"`
	ContactID uuid.UUID `json:"contact_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository define o contrato de acesso a dados para Lead.
type Repository interface {
	FindByID(id uuid.UUID) (*Lead, error)
	FindAll(ownerID uuid.UUID, filters Filters) ([]*Lead, int64, error)
	FindByContact(contactID uuid.UUID) ([]*Lead, error)
	Save(lead *Lead) (*Lead, error)
	Update(lead *Lead) (*Lead, error)
	Delete(id uuid.UUID) error
}

// Filters encapsula os parâmetros de pesquisa para leads.
// Embebe pagination.Base para herdar Page, Limit, SortBy, SortDir, Offset().
type Filters struct {
	pagination.Base
	Status Status
}

func (f *Filters) SetDefaults() {
	f.Normalize("created_at")
}
