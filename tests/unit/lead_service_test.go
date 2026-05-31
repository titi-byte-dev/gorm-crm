package unit_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/lead"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

var _ lead.Repository = (*leadRepoMem)(nil)

type leadRepoMem struct {
	data map[uuid.UUID]*lead.Lead
}

func newLeadRepo() *leadRepoMem {
	return &leadRepoMem{data: make(map[uuid.UUID]*lead.Lead)}
}

func (r *leadRepoMem) Save(l *lead.Lead) (*lead.Lead, error) {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	r.data[l.ID] = l
	return l, nil
}

func (r *leadRepoMem) FindByID(id uuid.UUID) (*lead.Lead, error) {
	l, ok := r.data[id]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return l, nil
}

func (r *leadRepoMem) FindAll(_ uuid.UUID, _ lead.Filters) ([]*lead.Lead, int64, error) {
	var out []*lead.Lead
	for _, l := range r.data {
		out = append(out, l)
	}
	return out, int64(len(out)), nil
}

func (r *leadRepoMem) FindByContact(_ uuid.UUID) ([]*lead.Lead, error) { return nil, nil }

func (r *leadRepoMem) Update(l *lead.Lead) (*lead.Lead, error) {
	r.data[l.ID] = l
	return l, nil
}

func (r *leadRepoMem) Delete(id uuid.UUID) error {
	delete(r.data, id)
	return nil
}

func newLeadService() *lead.Service {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return lead.NewService(newLeadRepo(), events.New(10, log))
}

// ---

// TestLeadService_UpdateStatus usa table-driven tests — o idioma Go para cobrir
// múltiplos cenários num único loop, sem repetição de setup/assert.
func TestLeadService_UpdateStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		transitions []lead.Status // sequência de transições a executar
		wantErr     bool
		errIs       error
	}{
		{
			name:        "new → contacted: válido",
			transitions: []lead.Status{lead.StatusContacted},
			wantErr:     false,
		},
		{
			name:        "new → contacted → qualified: válido",
			transitions: []lead.Status{lead.StatusContacted, lead.StatusQualified},
			wantErr:     false,
		},
		{
			name:        "new → qualified: inválido (salto de estado)",
			transitions: []lead.Status{lead.StatusQualified},
			wantErr:     true,
			errIs:       sharederrors.ErrValidation,
		},
		{
			name:        "new → lost: válido",
			transitions: []lead.Status{lead.StatusLost},
			wantErr:     false,
		},
		{
			name:        "lost → new: inválido (estado final)",
			transitions: []lead.Status{lead.StatusLost, lead.StatusNew},
			wantErr:     true,
			errIs:       sharederrors.ErrValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newLeadService()
			ownerID := uuid.New()

			l, err := svc.Create(ownerID, lead.CreateLeadDTO{
				Title:     "Lead teste",
				ContactID: uuid.New(),
			})
			if err != nil {
				t.Fatalf("create lead: %v", err)
			}

			var lastErr error
			for _, status := range tt.transitions {
				_, lastErr = svc.UpdateStatus(l.ID, status)
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

func TestLeadService_Create_SetsStatusNew(t *testing.T) {
	t.Parallel()
	svc := newLeadService()

	l, err := svc.Create(uuid.New(), lead.CreateLeadDTO{
		Title:     "Pipeline inicial",
		ContactID: uuid.New(),
		Value:     1500.0,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if l.Status != lead.StatusNew {
		t.Errorf("status = %s, want %s", l.Status, lead.StatusNew)
	}
}
