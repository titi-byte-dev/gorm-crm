// Package decorator implementa o Decorator pattern para repositórios.
//
// Decorator pattern: adiciona comportamento a um objecto sem modificar a sua
// implementação — e sem que o caller saiba que existe um decorator.
// A assinatura é idêntica à da interface original.
package decorator

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
)

// ContactRepoLogger decora qualquer contact.Repository com logging de timing.
//
// Transparência: ContactRepoLogger implementa contact.Repository.
// O Service não sabe se está a falar com o repositório real ou com o decorator.
type ContactRepoLogger struct {
	inner  contact.Repository
	logger *slog.Logger
}

// NewContactRepoLogger envolve qualquer contact.Repository com logging.
// inner pode ser postgresRepository, um mock, ou outro decorator.
func NewContactRepoLogger(inner contact.Repository, logger *slog.Logger) contact.Repository {
	return &ContactRepoLogger{inner: inner, logger: logger}
}

// Compile-time: garante que ContactRepoLogger implementa contact.Repository.
var _ contact.Repository = (*ContactRepoLogger)(nil)

func (d *ContactRepoLogger) Save(c *contact.Contact) (*contact.Contact, error) {
	start := time.Now()
	result, err := d.inner.Save(c)
	d.log("Save", time.Since(start), err)
	return result, err
}

func (d *ContactRepoLogger) FindByID(id uuid.UUID) (*contact.Contact, error) {
	start := time.Now()
	result, err := d.inner.FindByID(id)
	d.log("FindByID", time.Since(start), err)
	return result, err
}

func (d *ContactRepoLogger) FindAll(ownerID uuid.UUID, filters contact.Filters) ([]*contact.Contact, int64, error) {
	start := time.Now()
	result, total, err := d.inner.FindAll(ownerID, filters)
	d.log("FindAll", time.Since(start), err, "total", total)
	return result, total, err
}

func (d *ContactRepoLogger) FindByEmail(email string) (*contact.Contact, error) {
	start := time.Now()
	result, err := d.inner.FindByEmail(email)
	d.log("FindByEmail", time.Since(start), err)
	return result, err
}

func (d *ContactRepoLogger) Update(c *contact.Contact) (*contact.Contact, error) {
	start := time.Now()
	result, err := d.inner.Update(c)
	d.log("Update", time.Since(start), err)
	return result, err
}

func (d *ContactRepoLogger) Delete(id uuid.UUID) error {
	start := time.Now()
	err := d.inner.Delete(id)
	d.log("Delete", time.Since(start), err)
	return err
}

func (d *ContactRepoLogger) log(op string, elapsed time.Duration, err error, extra ...any) {
	args := append([]any{"op", op, "elapsed_ms", elapsed.Milliseconds()}, extra...)
	if err != nil {
		args = append(args, "error", err)
		d.logger.Warn("contact repo", args...)
		return
	}
	d.logger.Debug("contact repo", args...)
}
