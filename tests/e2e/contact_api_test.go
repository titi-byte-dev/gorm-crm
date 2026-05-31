package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
)

// newTestApp monta uma app Fiber minima com as rotas de contacts.
//
// Por que fiber.Test() em vez de httptest.Server?
//   - Nao abre portas de rede — o teste e mais rapido e isolado
//   - fiber.Test() executa o handler no mesmo processo
//   - Ideal para testar a camada HTTP sem base de dados
//
// O middleware de autenticacao e substituido por injectOwner,
// que injeta um ownerID fixo no contexto Fiber — sem JWT, sem chaves.
func newTestApp(t *testing.T) (*fiber.App, uuid.UUID) {
	t.Helper()

	ownerID := uuid.New()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	bus := events.New(10, log)
	repo := newContactRepoForE2E()
	svc := contact.NewService(repo, bus)

	app := fiber.New(fiber.Config{
		ErrorHandler: sharederrors.Handler,
	})

	// Middleware de teste: injeta ownerID sem validar JWT
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("userID", ownerID)
		return c.Next()
	})

	contact.RegisterRoutes(app, svc)

	return app, ownerID
}

// ---

func TestContactAPI_Create_201(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)

	body := mustJSON(t, map[string]string{
		"name":  "Ana Ferreira",
		"email": fmt.Sprintf("ana+%s@exemplo.pt", uuid.New()),
	})

	resp := doRequest(t, app, http.MethodPost, "/contacts", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("status = %d, want 201", resp.StatusCode)
	}

	var created contact.Contact
	mustDecode(t, resp.Body, &created)
	if created.ID == uuid.Nil {
		t.Error("ID nao deve ser nil apos criacao")
	}
	if created.Name != "Ana Ferreira" {
		t.Errorf("name = %s, want Ana Ferreira", created.Name)
	}
}

func TestContactAPI_Create_422_MissingName(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)

	body := mustJSON(t, map[string]string{
		"email": "sem-nome@exemplo.pt",
		// name em falta — validate:"required"
	})

	resp := doRequest(t, app, http.MethodPost, "/contacts", body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", resp.StatusCode)
	}
}

func TestContactAPI_Create_409_DuplicateEmail(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)
	email := fmt.Sprintf("dup+%s@exemplo.pt", uuid.New())

	body := mustJSON(t, map[string]string{"name": "Primeiro", "email": email})
	resp1 := doRequest(t, app, http.MethodPost, "/contacts", body)
	resp1.Body.Close()

	resp2 := doRequest(t, app, http.MethodPost, "/contacts",
		mustJSON(t, map[string]string{"name": "Segundo", "email": email}))
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusConflict {
		t.Errorf("status = %d, want 409", resp2.StatusCode)
	}
}

func TestContactAPI_GetByID_200(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)

	// Criar primeiro
	created := createContact(t, app, "Carlos Costa", fmt.Sprintf("carlos+%s@exemplo.pt", uuid.New()))

	// Ler por ID
	resp := doRequest(t, app, http.MethodGet, "/contacts/"+created.ID.String(), nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestContactAPI_GetByID_404(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)

	resp := doRequest(t, app, http.MethodGet, "/contacts/"+uuid.New().String(), nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", resp.StatusCode)
	}
}

func TestContactAPI_Delete_204(t *testing.T) {
	t.Parallel()
	app, _ := newTestApp(t)

	created := createContact(t, app, "Para Apagar", fmt.Sprintf("apagar+%s@exemplo.pt", uuid.New()))

	resp := doRequest(t, app, http.MethodDelete, "/contacts/"+created.ID.String(), nil)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("status = %d, want 204", resp.StatusCode)
	}

	// Confirma que nao existe mais
	get := doRequest(t, app, http.MethodGet, "/contacts/"+created.ID.String(), nil)
	get.Body.Close()
	if get.StatusCode != http.StatusNotFound {
		t.Errorf("after delete: status = %d, want 404", get.StatusCode)
	}
}

// --- helpers ---

// contactRepoForE2E e o mesmo mock que os unit tests, reutilizado aqui.
// Em Go, o mesmo tipo pode ser declarado uma vez e importado; neste caso
// declaramos localmente para manter o pacote e2e independente.
var _ contact.Repository = (*contactRepoForE2E)(nil)

type contactRepoForE2E struct {
	data    map[uuid.UUID]*contact.Contact
	byEmail map[string]*contact.Contact
}

func newContactRepoForE2E() *contactRepoForE2E {
	return &contactRepoForE2E{
		data:    make(map[uuid.UUID]*contact.Contact),
		byEmail: make(map[string]*contact.Contact),
	}
}

func (r *contactRepoForE2E) Save(c *contact.Contact) (*contact.Contact, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	r.data[c.ID] = c
	r.byEmail[c.Email] = c
	return c, nil
}

func (r *contactRepoForE2E) FindByID(id uuid.UUID) (*contact.Contact, error) {
	c, ok := r.data[id]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}

func (r *contactRepoForE2E) FindAll(_ uuid.UUID, _ contact.Filters) ([]*contact.Contact, int64, error) {
	var out []*contact.Contact
	for _, c := range r.data {
		out = append(out, c)
	}
	return out, int64(len(out)), nil
}

func (r *contactRepoForE2E) FindByEmail(email string) (*contact.Contact, error) {
	c, ok := r.byEmail[email]
	if !ok {
		return nil, sharederrors.ErrNotFound
	}
	return c, nil
}

func (r *contactRepoForE2E) Update(c *contact.Contact) (*contact.Contact, error) {
	r.data[c.ID] = c
	return c, nil
}

func (r *contactRepoForE2E) Delete(id uuid.UUID) error {
	c, ok := r.data[id]
	if !ok {
		return sharederrors.ErrNotFound
	}
	delete(r.byEmail, c.Email)
	delete(r.data, id)
	return nil
}

func doRequest(t *testing.T, app *fiber.App, method, path string, body []byte) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}
	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func mustDecode(t *testing.T, r io.Reader, v any) {
	t.Helper()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatalf("decode: %v", err)
	}
}

func createContact(t *testing.T, app *fiber.App, name, email string) contact.Contact {
	t.Helper()
	body := mustJSON(t, map[string]string{"name": name, "email": email})
	resp := doRequest(t, app, http.MethodPost, "/contacts", body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("createContact: status = %d", resp.StatusCode)
	}
	var c contact.Contact
	mustDecode(t, resp.Body, &c)
	return c
}
