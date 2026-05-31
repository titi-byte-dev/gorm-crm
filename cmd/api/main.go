package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
)

func main() {
	app := fiber.New(fiber.Config{
		AppName:      "GoRM CRM v0.1.0",
		ErrorHandler: errors.Handler,
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} ${method} ${path} (${latency})\n",
	}))

	registerRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 GoRM a correr em http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

func registerRoutes(app *fiber.App) {
	app.Get("/health", healthHandler)

	// v1 API group — será expandido nos módulos seguintes
	v1 := app.Group("/api/v1")
	_ = v1
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "gorm-crm",
		"version": "0.1.0",
	})
}
