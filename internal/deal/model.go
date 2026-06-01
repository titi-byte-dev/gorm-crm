package deal

import (
	"time"

	"github.com/google/uuid"
)

// Stage representa a etapa de um negócio no pipeline de vendas.
type Stage string

const (
	StageProposal    Stage = "proposal"
	StageNegotiation Stage = "negotiation"
	StageWon         Stage = "won"
	StageLost        Stage = "lost"
)

func (s Stage) IsValid() bool {
	switch s {
	case StageProposal, StageNegotiation, StageWon, StageLost:
		return true
	}
	return false
}

func (s Stage) IsClosed() bool {
	return s == StageWon || s == StageLost
}

var dealTransitions = map[Stage]map[Stage]bool{
	StageProposal:    {StageNegotiation: true, StageLost: true},
	StageNegotiation: {StageWon: true, StageLost: true},
	StageWon:         {},
	StageLost:        {},
}

func (s Stage) CanTransitionTo(next Stage) bool {
	return dealTransitions[s][next]
}

// Deal representa um negócio em curso no CRM.
type Deal struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	Value     float64    `json:"value"`
	Stage     Stage      `json:"stage"`
	LeadID    *uuid.UUID `json:"lead_id,omitempty"` // pointer — pode ser nil (deal sem lead)
	ContactID uuid.UUID  `json:"contact_id"`
	OwnerID   uuid.UUID  `json:"owner_id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	ClosedAt  *time.Time `json:"closed_at,omitempty"` // pointer — nil até fechar
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Repository define o contrato de acesso a dados para Deal.
type Repository interface {
	FindByID(id uuid.UUID) (*Deal, error)
	FindAll(tenantID, ownerID uuid.UUID, isManager bool, filters Filters) ([]*Deal, int64, error)
	FindByContact(contactID uuid.UUID) ([]*Deal, error)
	Save(deal *Deal) (*Deal, error)
	Update(deal *Deal) (*Deal, error)
	Delete(id uuid.UUID) error
}

// Filters encapsula os parâmetros de pesquisa para deals.
type Filters struct {
	Stage   Stage
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
