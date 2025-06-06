package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"siyahsensei/wallet-service/api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/rs/zerolog/log"

	"siyahsensei/wallet-service/configs"
	"siyahsensei/wallet-service/domain/account"
	"siyahsensei/wallet-service/domain/asset"
	"siyahsensei/wallet-service/domain/definition"
	"siyahsensei/wallet-service/domain/user"
	"siyahsensei/wallet-service/infrastructure/configuration/auth"
	"siyahsensei/wallet-service/infrastructure/configuration/database"
	customLogger "siyahsensei/wallet-service/infrastructure/configuration/logger"
	"siyahsensei/wallet-service/infrastructure/persistence/accountrepo"
	"siyahsensei/wallet-service/infrastructure/persistence/assetrepo"
	"siyahsensei/wallet-service/infrastructure/persistence/definitionrepo"
	"siyahsensei/wallet-service/infrastructure/persistence/userrepo"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	customLogger.InitLogger(customLogger.Config{
		LogLevel: "debug",
		Pretty:   config.Environment == "development",
	})

	customLogger.Info("Starting wallet API server...", map[string]interface{}{
		"environment": config.Environment,
		"port":        config.ServerPort,
	})
	db, err := database.NewPostgresDB(database.PostgresConfig{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
		SSLMode:  "disable",
	})

	if err != nil {
		customLogger.Fatal("Failed to connect to database", err)
	}
	defer db.Close()

	userRepo := userrepo.NewPostgresRepository(db)
	userService := user.NewHandler(userRepo, config.JWTSecret, config.TokenExpiry)

	definitionRepo := definitionrepo.NewPostgresRepository(db)
	definitionService := definition.NewHandler(definitionRepo)

	accountRepo := accountrepo.NewPostgresRepository(db)
	accountService := account.NewHandler(accountRepo)

	assetRepo := assetrepo.NewPostgresRepository(db)
	assetService := asset.NewHandler(assetRepo)

	jwtMiddleware := auth.NewJWTMiddleware(config.JWTSecret)
	app := fiber.New(fiber.Config{
		AppName:               "Wallet API",
		DisableStartupMessage: true,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	authRoute := routes.NewAuthRoute(userService, jwtMiddleware)
	definitionHandler := routes.NewDefinitionHandler(definitionService)
	accountHandler := routes.NewAccountHandler(accountService)
	assetHandler := routes.NewAssetHandler(assetService)

	api := app.Group("/api")
	authRoute.RegisterRoutes(api.Group("/auth"), jwtMiddleware.Middleware())
	definitionHandler.RegisterRoutes(api)
	accountHandler.RegisterRoutes(api, jwtMiddleware.Middleware())
	assetHandler.RegisterRoutes(api, jwtMiddleware.Middleware())
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	go func() {
		if err := app.Listen(":" + config.ServerPort); err != nil {
			customLogger.Fatal("Failed to start server", err)
		}
	}()

	customLogger.Info(fmt.Sprintf("Server started on port %s", config.ServerPort))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	customLogger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		customLogger.Fatal("Server shutdown failed", err)
	}

	customLogger.Info("Server gracefully stopped")
}
