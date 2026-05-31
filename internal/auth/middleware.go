package auth

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
)

// Aliases para as chaves de contexto definidas em ctxutil.
const (
	ContextUserID = ctxutil.KeyUserID
	ContextRole   = ctxutil.KeyRole
)

// Protected é o middleware de autenticação JWT.
// Extrai o token do header Authorization: Bearer <token>,
// valida-o e injeta userID e role no contexto do request.
//
// Padrão middleware Go/Fiber:
//   func(c *fiber.Ctx) error → processa → c.Next() → próximo handler
//
// Qualquer route group que use Protected() só é acessível com token válido.
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return sharederrors.ErrUnauthorized
		}

		// Formato obrigatório: "Bearer <token>"
		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return sharederrors.ErrUnauthorized
		}

		claims, err := ValidateToken(parts[1])
		if err != nil {
			return sharederrors.ErrUnauthorized
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return sharederrors.ErrUnauthorized
		}

		// Injeta no contexto — acessível em qualquer handler downstream
		c.Locals(ContextUserID, userID)
		c.Locals(ContextRole, claims.Role)

		return c.Next()
	}
}

// RequireRole é o middleware de autorização RBAC.
// Deve ser usado DEPOIS de Protected() — assume que o contexto já tem role.
//
// RBAC (Role-Based Access Control): permissões atribuídas a roles,
// não a utilizadores individuais. Fácil de auditar e gerir.
//
// Hierarquia:  admin > manager > seller
// Um admin pode fazer tudo o que um manager pode, etc.
func RequireRole(roles ...user.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals(ContextRole).(user.Role)
		if !ok {
			return sharederrors.ErrForbidden
		}
		for _, allowed := range roles {
			if role == allowed || isHigherRole(role, allowed) {
				return c.Next()
			}
		}
		return sharederrors.ErrForbidden
	}
}

// isHigherRole verifica se r tem privilégios iguais ou superiores a target.
// admin > manager > seller
func isHigherRole(r, target user.Role) bool {
	hierarchy := map[user.Role]int{
		user.RoleSeller:  1,
		user.RoleManager: 2,
		user.RoleAdmin:   3,
	}
	return hierarchy[r] >= hierarchy[target]
}

// UserIDFromCtx extrai o userID do contexto — helper para os handlers.
func UserIDFromCtx(c *fiber.Ctx) (uuid.UUID, error) {
	id, ok := c.Locals(ContextUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	return id, nil
}
