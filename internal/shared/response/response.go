package response

import "github.com/gofiber/fiber/v2"

// Page é o envelope standard para listas paginadas.
// Todos os endpoints de listagem devolvem este formato.
type Page[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Pages int64 `json:"pages"`
}

// NewPage constrói um Page calculando o número total de páginas.
func NewPage[T any](data []T, total int64, page, limit int) Page[T] {
	pages := total / int64(limit)
	if total%int64(limit) > 0 {
		pages++
	}
	return Page[T]{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
		Pages: pages,
	}
}

// Created devolve 201 com o recurso criado.
func Created(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(data)
}

// OK devolve 200 com dados.
func OK(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(data)
}

// NoContent devolve 204 sem body.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
