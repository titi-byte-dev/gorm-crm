package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/validate"
)

type Handler struct{ svc *Service }

func RegisterRoutes(router fiber.Router, svc *Service) {
	h := &Handler{svc: svc}
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Get("/me", Protected(), h.Me)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var dto RegisterDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	u, err := h.svc.Register(dto)
	if err != nil {
		return err
	}
	// Devolve o utilizador SEM o password_hash (tag json:"-" no model)
	return response.Created(c, u)
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var dto LoginDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	tokens, err := h.svc.Login(dto)
	if err != nil {
		return err
	}
	return response.OK(c, tokens)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	tokens, err := h.svc.Refresh(body.RefreshToken)
	if err != nil {
		return err
	}
	return response.OK(c, tokens)
}

// Me devolve os dados do utilizador autenticado — GET /api/v1/auth/me
// Requer Protected() middleware — extrai userID do contexto JWT.
func (h *Handler) Me(c *fiber.Ctx) error {
	userID, err := UserIDFromCtx(c)
	if err != nil {
		return err
	}
	role := c.Locals(ContextRole)
	return response.OK(c, fiber.Map{
		"user_id": userID,
		"role":    role,
	})
}
