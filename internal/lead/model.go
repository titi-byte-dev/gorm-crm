package lead

import (
	"time"

	"github.com/google/uuid"
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

var leadTransitions = map[Status]map[Status]bool{
	StatusNew:       {StatusContacted: true, StatusLost: true},
	StatusContacted: {StatusQualified: true, StatusLost: true},
	StatusQualified: {StatusLost: true},
	StatusLost:      {},
}

func (s Status) CanTransitionTo(next Status) bool {
	return leadTransitions[s][next]
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
type Filters struct {
	Status  Status
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
