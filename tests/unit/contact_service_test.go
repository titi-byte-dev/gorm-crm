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

// contactRepoMem implementa contact.Repository em memória.
// var _ garante que o compilador valida a interface em tempo de build — não em runtime.
var _ contact.Repository = (*contactRepoMem)(nil)

type contactRepoMem struct {
	data   map[uuid.UUID]*contact.Contact
	byEmail map[string]*contact.Contact
}

func newContactRepo() *contactRepoMem {
	return &contactRepoMem{
		data:    make(map[uuid.UUID]*contact.Contact),
		byEmail: make(map[string]*contact.Contact),
	}
}

func (r *contactRepoMem) Save(c *contact.Contact) (*contact.Contact, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	r.data[c.ID] = c
	r.byEmail[c.Email] = c
	return c, nil
}

func (r *contactRepoMem) FindByID(id uuid.UUID) (*contact.Contact, error) {
	c, ok := r.data[id]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}

func (r *contactRepoMem) FindAll(_ uuid.UUID, _ contact.Filters) ([]*contact.Contact, int64, error) {
	var out []*contact.Contact
	for _, c := range r.data {
		out = append(out, c)
	}
	return out, int64(len(out)), nil
}

func (r *contactRepoMem) FindByEmail(email string) (*contact.Contact, error) {
	c, ok := r.byEmail[email]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}

func (r *contactRepoMem) Update(c *contact.Contact) (*contact.Contact, error) {
	r.data[c.ID] = c
	return c, nil
}

func (r *contactRepoMem) Delete(id uuid.UUID) error {
	c, ok := r.data[id]
	if !ok {
		return sharederrors.ErrNotFound
	}
	delete(r.byEmail, c.Email)
	delete(r.data, id)
	return nil
}

// newContactService cria Service com repositório em memória e bus descartável.
func newContactService() (*contact.Service, *contactRepoMem) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	bus := events.New(10, log)
	repo := newContactRepo()
	return contact.NewService(repo, bus), repo
}

// ---

func TestContactService_Create_Success(t *testing.T) {
	t.Parallel()
	svc, _ := newContactService()
	ownerID := uuid.New()

	c, err := svc.Create(ownerID, contact.CreateContactDTO{
		Name:  "Ana Ferreira",
		Email: "ana@exemplo.pt",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID == uuid.Nil {
		t.Error("expected non-nil ID after save")
	}
	if c.OwnerID != ownerID {
		t.Errorf("owner_id = %s, want %s", c.OwnerID, ownerID)
	}
}

func TestContactService_Create_DuplicateEmail(t *testing.T) {
	t.Parallel()
	svc, _ := newContactService()
	ownerID := uuid.New()
	dto := contact.CreateContactDTO{Name: "João Silva", Email: "joao@exemplo.pt"}

	if _, err := svc.Create(ownerID, dto); err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	// Segundo create com o mesmo email deve falhar com ErrConflict
	_, err := svc.Create(ownerID, dto)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	if !errors.Is(err, sharederrors.ErrConflict) {
		t.Errorf("expected ErrConflict, got: %v", err)
	}
}

func TestContactService_GetByID_NotFound(t *testing.T) {
	t.Parallel()
	svc, _ := newContactService()

	_, err := svc.GetByID(uuid.New())
	if !errors.Is(err, sharederrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestContactService_Update_PartialFields(t *testing.T) {
	t.Parallel()
	svc, _ := newContactService()
	ownerID := uuid.New()

	created, _ := svc.Create(ownerID, contact.CreateContactDTO{
		Name:    "Carlos Costa",
		Email:   "carlos@exemplo.pt",
		Company: "Antiga Empresa",
	})

	newCompany := "Nova Empresa"
	updated, err := svc.Update(created.ID, contact.UpdateContactDTO{Company: &newCompany})
	if err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if updated.Name != "Carlos Costa" {
		t.Errorf("name changed unexpectedly: %s", updated.Name)
	}
	if updated.Company != newCompany {
		t.Errorf("company = %s, want %s", updated.Company, newCompany)
	}
}

func TestContactService_Delete_ThenNotFound(t *testing.T) {
	t.Parallel()
	svc, _ := newContactService()
	ownerID := uuid.New()

	created, _ := svc.Create(ownerID, contact.CreateContactDTO{
		Name:  "Maria Sousa",
		Email: "maria@exemplo.pt",
	})

	if err := svc.Delete(created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := svc.GetByID(created.ID)
	if !errors.Is(err, sharederrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}
