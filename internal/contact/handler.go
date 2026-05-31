package contact

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
)

// Handler liga os endpoints HTTP ao Service.
// Não contém lógica de negócio — só parsing, validação e serialização.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes regista todas as rotas de contactos no grupo v1.
func RegisterRoutes(router fiber.Router, svc *Service) {
	h := NewHandler(svc)
	g := router.Group("/contacts")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	// ownerID virá do JWT no Módulo 06 — por agora usamos um UUID fixo
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	var dto CreateContactDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	contact, err := h.svc.Create(ownerID, dto)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(contact)
}

func (h *Handler) List(c *fiber.Ctx) error {
	ownerID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	filters := Filters{
		Search:  c.Query("search"),
		Company: c.Query("company"),
		Page:    c.QueryInt("page", 1),
		Limit:   c.QueryInt("limit", 20),
		SortBy:  c.Query("sort_by", "created_at"),
		SortDir: c.Query("sort_dir", "desc"),
	}

	contacts, total, err := h.svc.List(ownerID, filters)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"data":  contacts,
		"total": total,
		"page":  filters.Page,
		"limit": filters.Limit,
	})
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact id")
	}

	contact, err := h.svc.GetByID(id)
	if err != nil {
		return err
	}

	return c.JSON(contact)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact id")
	}

	var dto UpdateContactDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	contact, err := h.svc.Update(id, dto)
	if err != nil {
		return err
	}

	return c.JSON(contact)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact id")
	}

	if err := h.svc.Delete(id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// parseOwnerID será substituído por middleware JWT no Módulo 06.
func parseOwnerID(c *fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals("userID")
	if raw == nil {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	id, ok := raw.(uuid.UUID)
	if !ok {
		return uuid.Nil, sharederrors.ErrUnauthorized
	}
	return id, nil
}
