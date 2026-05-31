package ctxutil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
)

const (
	KeyUserID = "userID"
	KeyRole   = "userRole"
)

// OwnerID extrai o userID injetado pelo JWT middleware.
// Todos os handlers que precisam do utilizador autenticado usam esta função.
func OwnerID(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals(KeyUserID)
	if raw == nil {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	id, ok := raw.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	return id, nil
}
