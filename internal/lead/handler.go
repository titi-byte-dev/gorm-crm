package lead

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/validate"
	"github.com/titi-byte-dev/gorm-crm/pkg/pagination"
)

type Handler struct{ svc *Service }

func RegisterRoutes(router fiber.Router, svc *Service) {
	h := &Handler{svc: svc}
	g := router.Group("/leads")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Patch("/:id/status", h.UpdateStatus)
	g.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	ownerID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	var dto CreateLeadDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	lead, err := h.svc.Create(ownerID, dto)
	if err != nil {
		return err
	}
	return response.Created(c, lead)
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
	}
	if s := c.Query("status"); s != "" {
		filters.Status = Status(s)
	}
	leads, total, err := h.svc.List(ownerID, filters)
	if err != nil {
		return err
	}
	return response.OK(c, response.NewPage(leads, total, filters.Page, filters.Limit))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lead id")
	}
	lead, err := h.svc.GetByID(id)
	if err != nil {
		return err
	}
	return response.OK(c, lead)
}

func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lead id")
	}
	var body struct {
		Status Status `json:"status" validate:"required,oneof=new contacted qualified lost"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(body); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	lead, err := h.svc.UpdateStatus(id, body.Status)
	if err != nil {
		return err
	}
	return response.OK(c, lead)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid lead id")
	}
	if err := h.svc.Delete(id); err != nil {
		return err
	}
	return response.NoContent(c)
}
