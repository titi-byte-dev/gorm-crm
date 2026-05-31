package integration_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// newTestDB inicia um container PostgreSQL efémero e devolve um *gorm.DB pronto a usar.
//
// testcontainers-go gere o ciclo de vida do container:
//   - inicia antes do teste
//   - termina automaticamente via t.Cleanup
//
// O container usa a imagem oficial postgres:16-alpine — pequena e rápida.
func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	ctx := context.Background()

	ctr, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		if err := ctr.Terminate(ctx); err != nil {
			t.Logf("terminate container: %v", err)
		}
	})

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open gorm: %v", err)
	}

	// AutoMigrate cria a tabela contacts — mesma estrutura que a migracao de producao.
	// Em testes de integracao, AutoMigrate e conveniente; em producao, usariamos
	// ficheiros de migracao numerados (ex: golang-migrate).
	if err := db.AutoMigrate(&contactRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return db
}

// contactRecord redefinido localmente para AutoMigrate sem importar o pacote interno.
// Replica a estrutura de contact/repository_pg.go.
type contactRecord struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name    string    `gorm:"not null"`
	Email   string    `gorm:"uniqueIndex;not null"`
	Phone   string
	Company string
	Notes   string
	OwnerID uuid.UUID `gorm:"type:uuid;not null;index"`
}

func (contactRecord) TableName() string { return "contacts" }

// ---

func TestContactRepository_SaveAndFind(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)
	ownerID := uuid.New()

	c := &contact.Contact{
		Name:    "Ana Ferreira",
		Email:   fmt.Sprintf("ana+%s@exemplo.pt", uuid.New()),
		Company: "ACME",
		OwnerID: ownerID,
	}

	saved, err := repo.Save(c)
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	if saved.ID == uuid.Nil {
		t.Fatal("ID deve ser preenchido apos save")
	}

	found, err := repo.FindByID(saved.ID)
	if err != nil {
		t.Fatalf("findByID: %v", err)
	}
	if found.Name != c.Name {
		t.Errorf("name = %s, want %s", found.Name, c.Name)
	}
	if found.Company != c.Company {
		t.Errorf("company = %s, want %s", found.Company, c.Company)
	}
}

func TestContactRepository_FindByID_NotFound(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)

	_, err := repo.FindByID(uuid.New())
	if !errors.Is(err, sharederrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got: %v", err)
	}
}

func TestContactRepository_FindByEmail(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)
	ownerID := uuid.New()
	email := fmt.Sprintf("unique+%s@exemplo.pt", uuid.New())

	_, err := repo.Save(&contact.Contact{
		Name:    "Teste",
		Email:   email,
		OwnerID: ownerID,
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}

	found, err := repo.FindByEmail(email)
	if err != nil {
		t.Fatalf("findByEmail: %v", err)
	}
	if found.Email != email {
		t.Errorf("email = %s, want %s", found.Email, email)
	}
}

func TestContactRepository_Update(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)
	ownerID := uuid.New()

	saved, _ := repo.Save(&contact.Contact{
		Name:    "Antes",
		Email:   fmt.Sprintf("update+%s@exemplo.pt", uuid.New()),
		OwnerID: ownerID,
	})

	saved.Name = "Depois"
	saved.Company = "Nova Empresa"
	updated, err := repo.Update(saved)
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Depois" {
		t.Errorf("name = %s, want Depois", updated.Name)
	}
}

func TestContactRepository_Delete(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)
	ownerID := uuid.New()

	saved, _ := repo.Save(&contact.Contact{
		Name:    "Para Apagar",
		Email:   fmt.Sprintf("delete+%s@exemplo.pt", uuid.New()),
		OwnerID: ownerID,
	})

	if err := repo.Delete(saved.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err := repo.FindByID(saved.ID)
	if !errors.Is(err, sharederrors.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got: %v", err)
	}
}

func TestContactRepository_FindAll_Pagination(t *testing.T) {
	t.Parallel()
	db := newTestDB(t)
	repo := contact.NewPostgresRepository(db)
	ownerID := uuid.New()

	// Insere 5 contactos para o mesmo owner
	for i := range 5 {
		_, err := repo.Save(&contact.Contact{
			Name:    fmt.Sprintf("Contacto %d", i),
			Email:   fmt.Sprintf("c%d+%s@exemplo.pt", i, uuid.New()),
			OwnerID: ownerID,
		})
		if err != nil {
			t.Fatalf("save %d: %v", i, err)
		}
	}

	// Pagina 1, limite 3
	contacts, total, err := repo.FindAll(ownerID, contact.Filters{Page: 1, Limit: 3})
	if err != nil {
		t.Fatalf("findAll: %v", err)
	}
	if total != 5 {
		t.Errorf("total = %d, want 5", total)
	}
	if len(contacts) != 3 {
		t.Errorf("len(contacts) = %d, want 3", len(contacts))
	}
}
