package contact

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

// contactRecord: índice composto (owner_id, created_at) porque FindAll usa
// WHERE owner_id = ? ORDER BY created_at — o b-tree cobre predicado e ordenação.
// Company indexado para o filtro ILIKE — parcial, mas evita full table scan no prefix.
type contactRecord struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     string
	Company   string    `gorm:"index"`
	Notes     string
	OwnerID   uuid.UUID `gorm:"type:uuid;not null;index:idx_contacts_owner_created,priority:1"`
	CreatedAt int64     `gorm:"autoCreateTime:milli;index:idx_contacts_owner_created,priority:2"`
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

func (r *postgresRepository) FindByIDs(ids []uuid.UUID) ([]*Contact, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var records []contactRecord
	if err := r.db.Where("id IN ?", ids).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("find contacts by ids: %w", err)
	}
	contacts := make([]*Contact, len(records))
	for i, rec := range records {
		contacts[i] = recordToContact(rec)
	}
	return contacts, nil
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
		Order(filters.SortBy + " " + filters.SortDir).
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
		ID:      r.ID,
		Name:    r.Name,
		Email:   r.Email,
		Phone:   r.Phone,
		Company: r.Company,
		Notes:   r.Notes,
		OwnerID: r.OwnerID,
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
