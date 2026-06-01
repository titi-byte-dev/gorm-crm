package contact

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

var allowedSortColumns = map[string]bool{
	"created_at": true,
	"updated_at": true,
	"name":       true,
	"email":      true,
	"company":    true,
}

func safeOrder(col, dir string) string {
	if !allowedSortColumns[col] {
		col = "created_at"
	}
	if strings.ToUpper(dir) == "ASC" {
		return col + " ASC"
	}
	return col + " DESC"
}

// contactRecord é o modelo GORM — separado do domain model para não
// vazar detalhes de persistência para o resto da aplicação.
type contactRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"not null;uniqueIndex:idx_contacts_email_owner"`
	Phone     string
	Company   string
	Notes     string
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index;uniqueIndex:idx_contacts_email_owner"`
	CreatedAt int64     `gorm:"autoCreateTime:milli"`
	UpdatedAt int64     `gorm:"autoUpdateTime:milli"`
}

func (contactRecord) TableName() string { return "contacts" }

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*Contact, error) {
	var rec contactRecord
	err := r.db.Where("id = ?", id).First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contact %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find contact: %w", err)
	}
	return recordToContact(rec), nil
}

func (r *postgresRepository) FindByEmail(email string) (*Contact, error) {
	var rec contactRecord
	err := r.db.Where("email = ?", email).First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contact email %s: %w", email, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find contact by email: %w", err)
	}
	return recordToContact(rec), nil
}

func (r *postgresRepository) FindAll(ownerID uuid.UUID, filters Filters) ([]*Contact, int64, error) {
	filters.SetDefaults()

	query := r.db.Model(&contactRecord{}).Where("owner_id = ?", ownerID)

	if filters.Search != "" {
		like := "%" + filters.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", like, like)
	}
	if filters.Company != "" {
		query = query.Where("company ILIKE ?", "%"+filters.Company+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count contacts: %w", err)
	}

	var records []contactRecord
	err := query.
		Order(safeOrder(filters.SortBy, filters.SortDir)).
		Limit(filters.Limit).
		Offset(filters.Offset()).
		Find(&records).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list contacts: %w", err)
	}

	contacts := make([]*Contact, len(records))
	for i, rec := range records {
		contacts[i] = recordToContact(rec)
	}
	return contacts, total, nil
}

func (r *postgresRepository) Save(contact *Contact) (*Contact, error) {
	if contact.ID == uuid.Nil {
		contact.ID = uuid.New()
	}
	rec := contactToRecord(contact)
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save contact: %w", err)
	}
	return recordToContact(rec), nil
}

func (r *postgresRepository) Update(contact *Contact) (*Contact, error) {
	rec := contactToRecord(contact)
	result := r.db.Model(&rec).Where("id = ?", contact.ID).Updates(rec)
	if result.Error != nil {
		return nil, fmt.Errorf("update contact: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("contact %s: %w", contact.ID, sharederrors.ErrNotFound)
	}
	return contact, nil
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&contactRecord{})
	if result.Error != nil {
		return fmt.Errorf("delete contact: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("contact %s: %w", id, sharederrors.ErrNotFound)
	}
	return nil
}

func recordToContact(r contactRecord) *Contact {
	return &Contact{
		ID:        r.ID,
		Name:      r.Name,
		Email:     r.Email,
		Phone:     r.Phone,
		Company:   r.Company,
		Notes:     r.Notes,
		OwnerID:   r.OwnerID,
		CreatedAt: time.UnixMilli(r.CreatedAt),
		UpdatedAt: time.UnixMilli(r.UpdatedAt),
	}
}

func contactToRecord(c *Contact) contactRecord {
	return contactRecord{
		ID:      c.ID,
		Name:    c.Name,
		Email:   c.Email,
		Phone:   c.Phone,
		Company: c.Company,
		Notes:   c.Notes,
		OwnerID: c.OwnerID,
	}
}
