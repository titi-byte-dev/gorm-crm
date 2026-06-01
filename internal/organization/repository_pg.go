package organization

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

type orgRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (orgRecord) TableName() string { return "organizations" }

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*Organization, error) {
	var rec orgRecord
	if err := r.db.First(&rec, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("organization %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find organization: %w", err)
	}
	return &Organization{ID: rec.ID, Name: rec.Name, CreatedAt: rec.CreatedAt}, nil
}

func (r *postgresRepository) Save(org *Organization) (*Organization, error) {
	if org.ID == uuid.Nil {
		org.ID = uuid.New()
	}
	rec := orgRecord{ID: org.ID, Name: org.Name}
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save organization: %w", err)
	}
	org.CreatedAt = rec.CreatedAt
	return org, nil
}
