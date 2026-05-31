package lead

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/pkg/valueobject"
	"gorm.io/gorm"
)

type leadRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title     string    `gorm:"not null"`
	Value     float64
	Status    string    `gorm:"not null;default:'new'"`
	ContactID uuid.UUID `gorm:"type:uuid;not null;index"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (leadRecord) TableName() string { return "leads" }

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*Lead, error) {
	var rec leadRecord
	if err := r.db.First(&rec, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("lead %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find lead: %w", err)
	}
	return recordToLead(rec), nil
}

func (r *postgresRepository) FindAll(ownerID uuid.UUID, filters Filters) ([]*Lead, int64, error) {
	filters.SetDefaults()
	query := r.db.Model(&leadRecord{}).Where("owner_id = ?", ownerID)
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	var total int64
	query.Count(&total)
	var recs []leadRecord
	err := query.Order(filters.SortBy + " " + filters.SortDir).
		Limit(filters.Limit).Offset(filters.Offset()).
		Find(&recs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list leads: %w", err)
	}
	leads := make([]*Lead, len(recs))
	for i, r := range recs {
		leads[i] = recordToLead(r)
	}
	return leads, total, nil
}

func (r *postgresRepository) FindByContact(contactID uuid.UUID) ([]*Lead, error) {
	var recs []leadRecord
	if err := r.db.Where("contact_id = ?", contactID).Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("find leads by contact: %w", err)
	}
	leads := make([]*Lead, len(recs))
	for i, rec := range recs {
		leads[i] = recordToLead(rec)
	}
	return leads, nil
}

func (r *postgresRepository) Save(lead *Lead) (*Lead, error) {
	if lead.ID == uuid.Nil {
		lead.ID = uuid.New()
	}
	rec := leadToRecord(lead)
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save lead: %w", err)
	}
	return recordToLead(rec), nil
}

func (r *postgresRepository) Update(lead *Lead) (*Lead, error) {
	rec := leadToRecord(lead)
	result := r.db.Model(&rec).Where("id = ?", lead.ID).Updates(rec)
	if result.Error != nil {
		return nil, fmt.Errorf("update lead: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("lead %s: %w", lead.ID, sharederrors.ErrNotFound)
	}
	return lead, nil
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&leadRecord{})
	if result.Error != nil {
		return fmt.Errorf("delete lead: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("lead %s: %w", id, sharederrors.ErrNotFound)
	}
	return nil
}

func recordToLead(r leadRecord) *Lead {
	return &Lead{
		ID: r.ID, Title: r.Title, Value: valueobject.Money(r.Value),
		Status: Status(r.Status), ContactID: r.ContactID, OwnerID: r.OwnerID,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func leadToRecord(l *Lead) leadRecord {
	return leadRecord{
		ID: l.ID, Title: l.Title, Value: l.Value.Float64(),
		Status: string(l.Status), ContactID: l.ContactID, OwnerID: l.OwnerID,
	}
}
