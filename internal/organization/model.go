package organization

import (
	"time"

	"github.com/google/uuid"
)

// Organization é o tenant root — todos os recursos pertencem a uma organização.
type Organization struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Repository interface {
	FindByID(id uuid.UUID) (*Organization, error)
	Save(org *Organization) (*Organization, error)
}
