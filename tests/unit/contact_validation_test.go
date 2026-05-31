package unit_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

func newContactSvc() *contact.Service {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return contact.NewService(newDecoratorContactRepo(), events.New(10, log))
}

func TestContactChain_UniqueEmail(t *testing.T) {
	t.Parallel()
	svc := newContactSvc()
	ownerID := uuid.New()
	dto := contact.CreateContactDTO{Name: "Ana", Email: "ana@exemplo.pt"}

	if _, err := svc.Create(ownerID, dto); err != nil {
		t.Fatalf("first create: %v", err)
	}
	_, err := svc.Create(ownerID, dto)
	if !errors.Is(err, sharederrors.ErrConflict) {
		t.Errorf("expected ErrConflict, got: %v", err)
	}
}

func TestContactChain_BlockedDomain(t *testing.T) {
	t.Parallel()
	svc := newContactSvc()

	_, err := svc.Create(uuid.New(), contact.CreateContactDTO{
		Name:  "Spam",
		Email: "test@mailinator.com",
	})
	if !errors.Is(err, sharederrors.ErrValidation) {
		t.Errorf("expected ErrValidation for blocked domain, got: %v", err)
	}
}

// TestContactChain_CustomRules prova que a chain e substituivel --
// o core do Chain of Responsibility: o caller controla a sequencia de regras.
func TestContactChain_CustomRules(t *testing.T) {
	t.Parallel()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	// Service sem EmailDomainRule -- aceita mailinator
	svc := contact.NewService(
		newDecoratorContactRepo(),
		events.New(10, log),
		contact.UniqueEmailRule{}, // chain personalizada: so unicidade
	)

	_, err := svc.Create(uuid.New(), contact.CreateContactDTO{
		Name:  "Teste",
		Email: "test@mailinator.com",
	})
	if err != nil {
		t.Errorf("expected nil with custom chain, got: %v", err)
	}
}
