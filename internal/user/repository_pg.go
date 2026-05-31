package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"gorm.io/gorm"
)

type userRecord struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         string    `gorm:"not null;default:'seller'"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (userRecord) TableName() string { return "users" }

var _ Repository = (*postgresRepository)(nil)

type postgresRepository struct{ db *gorm.DB }

func NewPostgresRepository(db *gorm.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) FindByID(id uuid.UUID) (*User, error) {
	var rec userRecord
	if err := r.db.First(&rec, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %s: %w", id, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return recordToUser(rec), nil
}

func (r *postgresRepository) FindByEmail(email string) (*User, error) {
	var rec userRecord
	if err := r.db.First(&rec, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user %s: %w", email, sharederrors.ErrNotFound)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return recordToUser(rec), nil
}

func (r *postgresRepository) Save(u *User) (*User, error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	rec := userToRecord(u)
	if err := r.db.Create(&rec).Error; err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}
	return recordToUser(rec), nil
}

func (r *postgresRepository) Update(u *User) (*User, error) {
	rec := userToRecord(u)
	if err := r.db.Model(&rec).Where("id = ?", u.ID).Updates(rec).Error; err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func recordToUser(r userRecord) *User {
	return &User{
		ID: r.ID, Name: r.Name, Email: r.Email,
		PasswordHash: r.PasswordHash, Role: Role(r.Role),
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func userToRecord(u *User) userRecord {
	return userRecord{
		ID: u.ID, Name: u.Name, Email: u.Email,
		PasswordHash: u.PasswordHash, Role: string(u.Role),
	}
}
