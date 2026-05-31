package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/titi-byte-dev/gorm-crm/internal/activitylog"
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
	"go.mongodb.org/mongo-driver/mongo"
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

	// MongoDB — opcional: a app funciona sem Mongo (logs descartados silenciosamente)
	var mongoDB *mongo.Database
	mongoDB, err = database.NewMongo(database.MongoConfigFromEnv())
	if err != nil {
		log.Warn("mongodb unavailable — activity logging disabled", "error", err)
	} else {
		log.Info("mongodb connected")
		actSvc := activitylog.NewService(activitylog.NewMongoRepository(mongoDB), log)
		actSvc.RegisterHandlers(bus) // subscreve eventos no bus
	}

	app := fiber.New(fiber.Config{
		AppName:      "GoRM CRM v0.9.0",
		ErrorHandler: sharederrors.Handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())

	registerRoutes(app, db, mongoDB, bus, log)

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

func registerRoutes(app *fiber.App, db *gorm.DB, mongoDB *mongo.Database, bus *events.Bus, log *slog.Logger) {
	app.Get("/health", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		dbStatus := "ok"
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "degraded"
		}
		mongoStatus := "disabled"
		if mongoDB != nil {
			mongoStatus = "ok"
		}
		status := "ok"
		httpStatus := fiber.StatusOK
		if dbStatus != "ok" {
			status = "degraded"
			httpStatus = fiber.StatusServiceUnavailable
		}
		return c.Status(httpStatus).JSON(fiber.Map{
			"status":  status,
			"service": "gorm-crm",
			"version": "0.9.0",
			"checks": fiber.Map{
				"database": dbStatus,
				"mongodb":  mongoStatus,
			},
		})
	})

	v1 := app.Group("/api/v1")

	authSvc := auth.NewService(user.NewPostgresRepository(db))
	auth.RegisterRoutes(v1, authSvc)

	protected := v1.Use(auth.Protected())

	contact.RegisterRoutes(protected, contact.NewService(contact.NewPostgresRepository(db), bus))
	lead.RegisterRoutes(protected, lead.NewService(lead.NewPostgresRepository(db), bus))
	deal.RegisterRoutes(protected, deal.NewService(deal.NewPostgresRepository(db), bus))
	task.RegisterRoutes(protected, task.NewService(task.NewPostgresRepository(db), bus))

	if mongoDB != nil {
		actSvc := activitylog.NewService(activitylog.NewMongoRepository(mongoDB), log)
		activitylog.RegisterRoutes(protected, actSvc)
	}
}
