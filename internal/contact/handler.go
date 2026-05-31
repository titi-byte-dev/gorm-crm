package contact

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/validate"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
)

type Handler struct{ svc *Service }

func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

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
	ownerID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	var dto CreateContactDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	contact, err := h.svc.Create(ownerID, dto)
	if err != nil {
		return err
	}
	return response.Created(c, contact)
}

func (h *Handler) List(c *fiber.Ctx) error {
	ownerID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	filters := Filters{
		Base: pagination.Base{
			Page:    c.QueryInt("page", 1),
			Limit:   c.QueryInt("limit", 20),
			SortBy:  c.Query("sort_by", "created_at"),
			SortDir: c.Query("sort_dir", "desc"),
		},
		Search:  c.Query("search"),
		Company: c.Query("company"),
	}
	contacts, total, err := h.svc.List(ownerID, filters)
	if err != nil {
		return err
	}
	return response.OK(c, response.NewPage(contacts, total, filters.Page, filters.Limit))
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
	return response.OK(c, contact)
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
	return response.OK(c, contact)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid contact id")
	}
	if err := h.svc.Delete(id); err != nil {
		return err
	}
	return response.NoContent(c)
}
