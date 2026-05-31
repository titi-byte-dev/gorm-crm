package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/middleware"
	"github.com/titi-byte-dev/gorm-crm/pkg/logger"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	log := logger.New(env)

	// Event bus com buffer de 500 — expandido com handlers no Módulo 09+
	bus := events.New(500, log)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bus.Start(ctx)

	app := fiber.New(fiber.Config{
		AppName:      "GoRM CRM v0.2.0",
		ErrorHandler: errors.Handler,
		// ReadTimeout evita slow-loris attacks
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())

	registerRoutes(app, bus)

	// Graceful shutdown — espera até 10s pelos requests em curso
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("shutting down server...")
		cancel() // sinaliza o event bus para terminar
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Error("forced shutdown", "error", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("server starting", "port", port, "env", env)
	if err := app.Listen(":" + port); err != nil {
		log.Error("server error", "error", err)
		os.Exit(1)
	}
}

func registerRoutes(app *fiber.App, bus *events.Bus) {
	app.Get("/health", healthHandler)

	v1 := app.Group("/api/v1")
	_ = v1
	// Rotas adicionadas progressivamente nos módulos seguintes:
	// M03: contact.RegisterRoutes(v1, db)
	// M05: lead.RegisterRoutes(v1, db)
	// M06: auth.RegisterRoutes(v1, db)
}

func healthHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "gorm-crm",
		"version": "0.2.0",
	})
}
