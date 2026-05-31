package errors

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// Sentinel errors — usados em toda a aplicação com errors.Is()
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation error")
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// Handler é o error handler global do Fiber.
// Mapeia erros de domínio para códigos HTTP.
func Handler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, ErrNotFound):
		code = fiber.StatusNotFound
		msg = err.Error()
	case errors.Is(err, ErrUnauthorized):
		code = fiber.StatusUnauthorized
		msg = err.Error()
	case errors.Is(err, ErrForbidden):
		code = fiber.StatusForbidden
		msg = err.Error()
	case errors.Is(err, ErrConflict):
		code = fiber.StatusConflict
		msg = err.Error()
	case errors.Is(err, ErrValidation):
		code = fiber.StatusUnprocessableEntity
		msg = err.Error()
	default:
		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			code = fiberErr.Code
			msg = fiberErr.Message
		}
	}

	return c.Status(code).JSON(ErrorResponse{
		Error:   http.StatusText(code),
		Message: msg,
	})
}
