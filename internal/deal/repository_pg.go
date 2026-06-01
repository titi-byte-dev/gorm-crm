package deal

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

var allowedDealSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"title":      true,
	"value":      true,
	"stage":      true,
	"closed_at":  true,
}

func safeDealOrder(col, dir string) string {
	if !allowedDealSortColumns[col] {
		col = "created_at"
	}
	if strings.ToUpper(dir) == "ASC" {
		return col + " ASC"
	}
	return col + " DESC"
}

type dealRecord struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Title     string     `gorm:"not null"`
	Value     float64
	Stage     string     `gorm:"not null;default:'proposal'"`
	LeadID    *uuid.UUID `gorm:"type:uuid;index"`
	ContactID uuid.UUID  `gorm:"type:uuid;not null;index"`
	OwnerID   uuid.UUID  `gorm:"type:uuid;not null;index"`
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index"`
	ClosedAt  *time.Time
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}

func (dealRecord) TableName() string { return "deals" }

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*Deal, error) {
	var rec dealRecord
	if err := r.db.First(&rec, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("deal %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find deal: %w", err)
	}
	return recordToDeal(rec), nil
}

func (r *postgresRepository) FindAll(tenantID, ownerID uuid.UUID, isManager bool, filters Filters) ([]*Deal, int64, error) {
	filters.SetDefaults()
	query := r.db.Model(&dealRecord{}).Where("tenant_id = ?", tenantID)
	if !isManager {
		query = query.Where("owner_id = ?", ownerID)
	}
	if filters.Stage != "" {
		query = query.Where("stage = ?", filters.Stage)
	}
	var total int64
	query.Count(&total)
	var recs []dealRecord
	err := query.Order(safeDealOrder(filters.SortBy, filters.SortDir)).
		Limit(filters.Limit).Offset((filters.Page - 1) * filters.Limit).
		Find(&recs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list deals: %w", err)
	}
	deals := make([]*Deal, len(recs))
	for i, rec := range recs {
		deals[i] = recordToDeal(rec)
	}
	return deals, total, nil
}

func (r *postgresRepository) FindByContact(contactID uuid.UUID) ([]*Deal, error) {
	var recs []dealRecord
	if err := r.db.Where("contact_id = ?", contactID).Find(&recs).Error; err != nil {
		return nil, fmt.Errorf("find deals by contact: %w", err)
	}
	deals := make([]*Deal, len(recs))
	for i, rec := range recs {
		deals[i] = recordToDeal(rec)
	}
	return deals, nil
}

func (r *postgresRepository) Save(deal *Deal) (*Deal, error) {
	if deal.ID == uuid.Nil {
		deal.ID = uuid.New()
	}
	rec := dealToRecord(deal)
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save deal: %w", err)
	}
	return recordToDeal(rec), nil
}

func (r *postgresRepository) Update(deal *Deal) (*Deal, error) {
	rec := dealToRecord(deal)
	result := r.db.Model(&rec).Where("id = ?", deal.ID).Updates(rec)
	if result.Error != nil {
		return nil, fmt.Errorf("update deal: %w", result.Error)
	}
	return deal, nil
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&dealRecord{})
	if result.Error != nil {
		return fmt.Errorf("delete deal: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("deal %s: %w", id, sharederrors.ErrNotFound)
	}
	return nil
}

func recordToDeal(r dealRecord) *Deal {
	return &Deal{
		ID: r.ID, Title: r.Title, Value: r.Value, Stage: Stage(r.Stage),
		LeadID: r.LeadID, ContactID: r.ContactID, OwnerID: r.OwnerID, TenantID: r.TenantID,
		ClosedAt: r.ClosedAt, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func dealToRecord(d *Deal) dealRecord {
	return dealRecord{
		ID: d.ID, Title: d.Title, Value: d.Value, Stage: string(d.Stage),
		LeadID: d.LeadID, ContactID: d.ContactID, OwnerID: d.OwnerID, TenantID: d.TenantID,
		ClosedAt: d.ClosedAt,
	}
}
