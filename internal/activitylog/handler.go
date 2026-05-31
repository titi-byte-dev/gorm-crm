package activitylog

import (
	"github.com/gofiber/fiber/v2"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/ctxutil"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/response"
)

type Handler struct{ svc *Service }

func RegisterRoutes(router fiber.Router, svc *Service) {
	h := &Handler{svc: svc}
	logs := router.Group("/activity")
	logs.Get("/me", h.MyActivity)
	logs.Get("/:entity_type/:entity_id", h.EntityActivity)
}

// MyActivity devolve o histórico de actividade do utilizador autenticado.
func (h *Handler) MyActivity(c *fiber.Ctx) error {
	userID, err := ctxutil.OwnerID(c)
	if err != nil {
		return err
	}
	limit := c.QueryInt("limit", 50)
	logs, err := h.svc.GetByUser(userID.String(), limit)
	if err != nil {
		return err
	}
	return response.OK(c, logs)
}

// EntityActivity devolve o histórico de uma entidade específica.
// GET /activity/contact/<uuid> → todos os eventos deste contacto
// GET /activity/deal/<uuid>    → todos os eventos deste deal
func (h *Handler) EntityActivity(c *fiber.Ctx) error {
	entityType := EntityType(c.Params("entity_type"))
	entityID := c.Params("entity_id")
	limit := c.QueryInt("limit", 50)

	logs, err := h.svc.GetByEntity(entityType, entityID, limit)
	if err != nil {
		return err
	}
	return response.OK(c, logs)
}
