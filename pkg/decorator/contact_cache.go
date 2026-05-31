package decorator

import (
	"time"

	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	"github.com/titi-byte-dev/gorm-crm/pkg/cache"
)

// CachingContactRepo é um Decorator que adiciona cache TTL ao contact.Repository.
// FindByID serve da cache; Save/Update/Delete invalidam a entrada afectada.
// FindByIDs, FindAll, FindByEmail delegam sempre ao inner (resultados não cacheados).
type CachingContactRepo struct {
	inner contact.Repository
	cache *cache.TTL[uuid.UUID, *contact.Contact]
}

var _ contact.Repository = (*CachingContactRepo)(nil)

func NewCachingContactRepo(inner contact.Repository, ttl time.Duration) contact.Repository {
	return &CachingContactRepo{inner: inner, cache: cache.New[uuid.UUID, *contact.Contact](ttl)}
}

func (r *CachingContactRepo) FindByID(id uuid.UUID) (*contact.Contact, error) {
	if cached, ok := r.cache.Get(id); ok {
		return cached, nil
	}
	c, err := r.inner.FindByID(id)
	if err == nil {
		r.cache.Set(id, c)
	}
	return c, err
}

func (r *CachingContactRepo) FindByIDs(ids []uuid.UUID) ([]*contact.Contact, error) {
	return r.inner.FindByIDs(ids)
}

func (r *CachingContactRepo) FindAll(ownerID uuid.UUID, filters contact.Filters) ([]*contact.Contact, int64, error) {
	return r.inner.FindAll(ownerID, filters)
}

func (r *CachingContactRepo) FindByEmail(email string) (*contact.Contact, error) {
	return r.inner.FindByEmail(email)
}

func (r *CachingContactRepo) Save(c *contact.Contact) (*contact.Contact, error) {
	saved, err := r.inner.Save(c)
	if err == nil {
		r.cache.Set(saved.ID, saved)
	}
	return saved, err
}

func (r *CachingContactRepo) Update(c *contact.Contact) (*contact.Contact, error) {
	updated, err := r.inner.Update(c)
	if err == nil {
		r.cache.Set(updated.ID, updated)
	}
	return updated, err
}

func (r *CachingContactRepo) Delete(id uuid.UUID) error {
	err := r.inner.Delete(id)
	if err == nil {
		r.cache.Delete(id)
	}
	return err
}
