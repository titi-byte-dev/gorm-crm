package ctxutil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

const (
	KeyUserID = "userID"
	KeyRole   = "userRole"
	KeyOrgID  = "orgID"
)

// RequestCtx encapsula o contexto do utilizador autenticado.
// Passado para os services em vez de parâmetros individuais.
type RequestCtx struct {
	UserID   uuid.UUID
	TenantID uuid.UUID
	Role     user.Role
}

// IsManager reporta se o utilizador pode ver dados de toda a organização.
func (r RequestCtx) IsManager() bool {
	return r.Role == user.RoleManager || r.Role == user.RoleAdmin
}

// FromFiber extrai o RequestCtx completo do contexto do request.
func FromFiber(c *fiber.Ctx) (RequestCtx, error) {
	userID, err := OwnerID(c)
	if err != nil {
		return RequestCtx{}, err
	}
	tenantID, err := TenantID(c)
	if err != nil {
		return RequestCtx{}, err
	}
	return RequestCtx{
		UserID:   userID,
		TenantID: tenantID,
		Role:     UserRole(c),
	}, nil
}

// OwnerID extrai o userID injetado pelo JWT middleware.
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

// TenantID extrai o orgID injetado pelo JWT middleware.
func TenantID(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals(KeyOrgID)
	if raw == nil {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	id, ok := raw.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	return id, nil
}

// UserRole extrai o role do utilizador autenticado.
func UserRole(c *fiber.Ctx) user.Role {
	role, _ := c.Locals(KeyRole).(user.Role)
	return role
}
