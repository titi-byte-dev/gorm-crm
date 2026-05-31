package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS configura Cross-Origin Resource Sharing.
// Em desenvolvimento permite qualquer origem; em produção só as origens listadas.
func CORS() fiber.Handler {
	allowedOrigins := os.Getenv("CORS_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*" // desenvolvimento
	}

	return cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodOptions,
		}, ","),
		AllowHeaders: "Origin, Content-Type, Authorization",
		MaxAge:       86400, // preflight cache: 24h
	})
}
