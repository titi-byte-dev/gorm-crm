package unit_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/deal"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

var _ deal.Repository = (*dealRepoMem)(nil)

type dealRepoMem struct {
	data map[uuid.UUID]*deal.Deal
}

func newDealRepo() *dealRepoMem {
	return &dealRepoMem{data: make(map[uuid.UUID]*deal.Deal)}
}

func (r *dealRepoMem) Save(d *deal.Deal) (*deal.Deal, error) {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	r.data[d.ID] = d
	return d, nil
}

func (r *dealRepoMem) FindByID(id uuid.UUID) (*deal.Deal, error) {
	d, ok := r.data[id]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return d, nil
}

func (r *dealRepoMem) FindAll(_ uuid.UUID, _ deal.Filters) ([]*deal.Deal, int64, error) {
	var out []*deal.Deal
	for _, d := range r.data {
		out = append(out, d)
	}
	return out, int64(len(out)), nil
}

func (r *dealRepoMem) FindByContact(_ uuid.UUID) ([]*deal.Deal, error) { return nil, nil }

func (r *dealRepoMem) Update(d *deal.Deal) (*deal.Deal, error) {
	r.data[d.ID] = d
	return d, nil
}

func (r *dealRepoMem) Delete(id uuid.UUID) error {
	delete(r.data, id)
	return nil
}

func newDealService() *deal.Service {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return deal.NewService(newDealRepo(), events.New(10, log))
}

// ---

// TestDealService_MoveStage verifica as transições de etapa do deal pipeline.
// Cada test case é independente — cria o seu próprio deal para evitar estado partilhado.
func TestDealService_MoveStage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		transitions []deal.Stage
		wantErr     bool
		errIs       error
	}{
		{
			name:        "proposal → negotiation: válido",
			transitions: []deal.Stage{deal.StageNegotiation},
			wantErr:     false,
		},
		{
			name:        "proposal → negotiation → won: válido",
			transitions: []deal.Stage{deal.StageNegotiation, deal.StageWon},
			wantErr:     false,
		},
		{
			name:        "proposal → negotiation → lost: válido",
			transitions: []deal.Stage{deal.StageNegotiation, deal.StageLost},
			wantErr:     false,
		},
		{
			name:        "proposal → won: inválido (salto de etapa)",
			transitions: []deal.Stage{deal.StageWon},
			wantErr:     true,
			errIs:       sharederrors.ErrValidation,
		},
		{
			name:        "won → lost: inválido (etapa final)",
			transitions: []deal.Stage{deal.StageNegotiation, deal.StageWon, deal.StageLost},
			wantErr:     true,
			errIs:       sharederrors.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newDealService()
			ownerID := uuid.New()

			d, err := svc.Create(ownerID, deal.CreateDealDTO{
				Title:     "Deal teste",
				ContactID: uuid.New(),
				Value:     5000.0,
			})
			if err != nil {
				t.Fatalf("create deal: %v", err)
			}

			var lastErr error
			for _, stage := range tt.transitions {
				_, lastErr = svc.MoveStage(d.ID, stage)
				if lastErr != nil {
					break
				}
			}

			if tt.wantErr {
				if lastErr == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errIs != nil && !errors.Is(lastErr, tt.errIs) {
					t.Errorf("expected %v, got: %v", tt.errIs, lastErr)
				}
			} else {
				if lastErr != nil {
					t.Errorf("unexpected error: %v", lastErr)
				}
			}
		})
	}
}

// TestDealService_MoveStage_ClosedAt verifica que ClosedAt é preenchido ao fechar.
// Testa um efeito colateral não-óbvio — exatamente o tipo de caso que os testes capturam.
func TestDealService_MoveStage_ClosedAt(t *testing.T) {
	t.Parallel()
	svc := newDealService()
	ownerID := uuid.New()

	d, _ := svc.Create(ownerID, deal.CreateDealDTO{
		Title:     "Fecho rápido",
		ContactID: uuid.New(),
	})
	if d.ClosedAt != nil {
		t.Error("ClosedAt deve ser nil antes de fechar")
	}

	// proposal → negotiation
	d, err := svc.MoveStage(d.ID, deal.StageNegotiation)
	if err != nil {
		t.Fatalf("move to negotiation: %v", err)
	}

	// negotiation → won
	d, err = svc.MoveStage(d.ID, deal.StageWon)
	if err != nil {
		t.Fatalf("move to won: %v", err)
	}

	if d.ClosedAt == nil {
		t.Error("ClosedAt deve ser preenchido ao fechar")
	}
}

// TestDealService_Create_DefaultStage verifica que deals começam em Proposal.
func TestDealService_Create_DefaultStage(t *testing.T) {
	t.Parallel()
	svc := newDealService()

	d, err := svc.Create(uuid.New(), deal.CreateDealDTO{
		Title:     "Proposta inicial",
		ContactID: uuid.New(),
		Value:     10000.0,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if d.Stage != deal.StageProposal {
		t.Errorf("stage = %s, want %s", d.Stage, deal.StageProposal)
	}
}
