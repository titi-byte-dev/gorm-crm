package deal

import (
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
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

func (s Stage) CanTransitionTo(next Stage) bool {
	transitions := map[Stage][]Stage{
		StageProposal:    {StageNegotiation, StageLost},
		StageNegotiation: {StageWon, StageLost},
		StageWon:         {},
		StageLost:        {},
	}
	for _, allowed := range transitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
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
	ClosedAt  *time.Time `json:"closed_at,omitempty"` // pointer — nil até fechar
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Reader define operacoes de leitura sobre Deal.
type Reader interface {
	FindByID(id uuid.UUID) (*Deal, error)
	FindAll(ownerID uuid.UUID, filters Filters) ([]*Deal, int64, error)
	FindByContact(contactID uuid.UUID) ([]*Deal, error)
}

// Writer define operacoes de escrita sobre Deal.
type Writer interface {
	Save(deal *Deal) (*Deal, error)
	Update(deal *Deal) (*Deal, error)
	Delete(id uuid.UUID) error
}

// Repository e a composicao de Reader e Writer.
type Repository interface {
	Reader
	Writer
}

// Filters encapsula os parâmetros de pesquisa para deals.
// Embebe pagination.Base para herdar Page, Limit, SortBy, SortDir, Offset().
type Filters struct {
	pagination.Base
	Stage Stage
}

func (f *Filters) SetDefaults() {
	f.Normalize("created_at")
}
