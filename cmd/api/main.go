package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/titi-byte-dev/gorm-crm/internal/auth"
	"github.com/titi-byte-dev/gorm-crm/internal/contact"
	"github.com/titi-byte-dev/gorm-crm/internal/deal"
	"github.com/titi-byte-dev/gorm-crm/internal/lead"
	"github.com/titi-byte-dev/gorm-crm/internal/task"
	sharederrors "github.com/titi-byte-dev/gorm-crm/internal/shared/errors"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/events"
	"github.com/titi-byte-dev/gorm-crm/internal/shared/middleware"
	"github.com/titi-byte-dev/gorm-crm/internal/user"
	"github.com/titi-byte-dev/gorm-crm/pkg/database"
	"github.com/titi-byte-dev/gorm-crm/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	log := logger.New(env)

	db, err := database.New(database.ConfigFromEnv(), env)
	if err != nil {
		log.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	log.Info("database connected")

	bus := events.New(500, log)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bus.Start(ctx)

	app := fiber.New(fiber.Config{
		AppName:      "GoRM CRM v0.7.0",
		ErrorHandler: sharederrors.Handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	registerRoutes(app, db, bus)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info("shutting down server...")
		cancel()
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

func registerRoutes(app *fiber.App, db *gorm.DB, bus *events.Bus) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "gorm-crm", "version": "0.7.0"})
	})

	v1 := app.Group("/api/v1")

	// Auth — rotas públicas (não requerem token)
	authSvc := auth.NewService(user.NewPostgresRepository(db))
	auth.RegisterRoutes(v1, authSvc)

	// Rotas protegidas — requerem JWT válido
	// auth.Protected() injeta userID e role no contexto de cada request
	protected := v1.Use(auth.Protected())

	contact.RegisterRoutes(protected, contact.NewService(contact.NewPostgresRepository(db), bus))
	lead.RegisterRoutes(protected, lead.NewService(lead.NewPostgresRepository(db), bus))
	deal.RegisterRoutes(protected, deal.NewService(deal.NewPostgresRepository(db), bus))
	task.RegisterRoutes(protected, task.NewService(task.NewPostgresRepository(db), bus))
}
