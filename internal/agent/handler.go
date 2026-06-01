package agent

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/validate"
)

type Handler struct{ svc *Service }

func RegisterRoutes(router fiber.Router, svc *Service) {
	h := &Handler{svc: svc}
	g := router.Group("/agents")
	g.Post("/run", h.Run)
	g.Post("/runs/:id/approve", h.Approve)
	g.Get("/runs", h.ListByEntity)
}

func (h *Handler) Run(c *fiber.Ctx) error {
	rctx, err := ctxutil.FromFiber(c)
	if err != nil {
		return err
	}
	var dto RunDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	run, err := h.svc.Run(rctx, dto)
	if err != nil {
		return err
	}
	return response.Created(c, run)
}

func (h *Handler) Approve(c *fiber.Ctx) error {
	rctx, err := ctxutil.FromFiber(c)
	if err != nil {
		return err
	}
	runID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid run id")
	}
	var dto ApproveDTO
	if err := c.BodyParser(&dto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if result := validate.Check(dto); result != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(result)
	}
	run, err := h.svc.ApproveActions(runID, rctx, dto.ActionIndices)
	if err != nil {
		return err
	}
	return response.OK(c, run)
}

func (h *Handler) ListByEntity(c *fiber.Ctx) error {
	rctx, err := ctxutil.FromFiber(c)
	if err != nil {
		return err
	}
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")
	if entityType == "" || entityIDStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "entity_type and entity_id are required")
	}
	entityID, err := uuid.Parse(entityIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid entity_id")
	}
	runs, err := h.svc.GetRunsByEntity(rctx, entityType, entityID)
	if err != nil {
		return err
	}
	return response.OK(c, runs)
}
