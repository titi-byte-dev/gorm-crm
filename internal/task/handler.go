package task

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/validate"
)

// Handler — só HTTP. Zero lógica de negócio, zero SQL.
type Handler struct{ svc *Service }

func RegisterRoutes(router fiber.Router, svc *Service) {
	h := &Handler{svc: svc}
	g := router.Group("/tasks")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/overdue", h.Overdue) // antes de /:id para não ser capturado como param
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Patch("/:id/status", h.UpdateStatus)
	g.Delete("/:id", h.Delete)
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var dto CreateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	task, err := h.svc.Create(dto)
	if err != nil {
		return err
	}
	return response.Created(c, task)
}

func (h *Handler) List(c *fiber.Ctx) error {
	assignedTo, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	filters := Filters{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 20),
	}
	if s := c.Query("status"); s != "" {
		filters.Status = Status(s)
	}
	if p := c.Query("priority"); p != "" {
		filters.Priority = Priority(p)
	}
	tasks, total, err := h.svc.List(assignedTo, filters)
	if err != nil {
		return err
	}
	return response.OK(c, response.NewPage(tasks, total, filters.Page, filters.Limit))
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	requesterID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}
	task, err := h.svc.GetByID(id, requesterID)
	if err != nil {
		return err
	}
	return response.OK(c, task)
}

func (h *Handler) Update(c *fiber.Ctx) error {
	requesterID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}
	var dto UpdateTaskDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	task, err := h.svc.Update(id, requesterID, dto)
	if err != nil {
		return err
	}
	return response.OK(c, task)
}

func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	requesterID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}
	var body struct {
		Status Status `json:"status" validate:"required,oneof=todo in_progress done cancelled"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(body); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	task, err := h.svc.UpdateStatus(id, requesterID, body.Status)
	if err != nil {
		return err
	}
	return response.OK(c, task)
}

func (h *Handler) Overdue(c *fiber.Ctx) error {
	tasks, err := h.svc.GetOverdue()
	if err != nil {
		return err
	}
	return response.OK(c, tasks)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	requesterID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid task id")
	}
	if err := h.svc.Delete(id, requesterID); err != nil {
		return err
	}
	return response.NoContent(c)
}
