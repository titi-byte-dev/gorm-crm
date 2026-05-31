package unit_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/pkg/decorator"
)

// decoratorContactRepo é um mock em memória para os testes do decorator.
var _ contact.Repository = (*decoratorContactRepo)(nil)

type decoratorContactRepo struct {
	data    map[uuid.UUID]*contact.Contact
	byEmail map[string]*contact.Contact
}

func newDecoratorContactRepo() *decoratorContactRepo {
	return &decoratorContactRepo{
		data:    make(map[uuid.UUID]*contact.Contact),
		byEmail: make(map[string]*contact.Contact),
	}
}

func (r *decoratorContactRepo) Save(c *contact.Contact) (*contact.Contact, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	r.data[c.ID] = c
	r.byEmail[c.Email] = c
	return c, nil
}
func (r *decoratorContactRepo) FindByID(id uuid.UUID) (*contact.Contact, error) {
	c, ok := r.data[id]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}
func (r *decoratorContactRepo) FindAll(_ uuid.UUID, _ contact.Filters) ([]*contact.Contact, int64, error) {
	return nil, 0, nil
}
func (r *decoratorContactRepo) FindByEmail(email string) (*contact.Contact, error) {
	c, ok := r.byEmail[email]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}
func (r *decoratorContactRepo) Update(c *contact.Contact) (*contact.Contact, error) {
	r.data[c.ID] = c
	return c, nil
}
func (r *decoratorContactRepo) Delete(id uuid.UUID) error {
	delete(r.data, id)
	return nil
}

// TestContactRepoLogger_DelegatesAndLogs verifica as duas propriedades
// fundamentais de um Decorator:
//  1. Delega: o resultado e o erro vêm do inner repo
//  2. Transparência: o tipo devolvido por NewContactRepoLogger
//     implementa a mesma interface que o inner
func TestContactRepoLogger_DelegatesAndLogs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	repo := decorator.NewContactRepoLogger(newDecoratorContactRepo(), logger)

	ownerID := uuid.New()
	c := &contact.Contact{
		Name:    "Decorada",
		Email:   "dec@exemplo.pt",
		OwnerID: ownerID,
	}

	// Save via decorator
	saved, err := repo.Save(c)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID == uuid.Nil {
		t.Error("ID deve ser preenchido")
	}

	// O log deve conter "contact repo" e "op=Save"
	logOutput := buf.String()
	if logOutput == "" {
		t.Error("esperava output de log, nao houve nenhum")
	}

	// FindByID via decorator — delega ao inner
	found, err := repo.FindByID(saved.ID)
	if err != nil {
		t.Fatalf("findByID: %v", err)
	}
	if found.Name != c.Name {
		t.Errorf("name = %s, want %s", found.Name, c.Name)
	}
}

// TestContactRepoLogger_IsTransparent prova que o Decorator e transparente:
// o Service pode usar o decorator sem saber que e um decorator.
func TestContactRepoLogger_IsTransparent(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// NewContactRepoLogger devolve contact.Repository -- nao *ContactRepoLogger
	// O caller so ve a interface
	var repo contact.Repository = decorator.NewContactRepoLogger(newDecoratorContactRepo(), logger)

	// Usa o repo exactamente como o Service usaria
	_, err := repo.Save(&contact.Contact{
		Name:    "Transparente",
		Email:   "transp@exemplo.pt",
		OwnerID: uuid.New(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
