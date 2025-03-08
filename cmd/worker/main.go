package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"siyahsensei/wallet-service/configs"
	"siyahsensei/wallet-service/domain/account"
	"siyahsensei/wallet-service/domain/asset"
	"siyahsensei/wallet-service/infrastructure/persistence/accountrepo"
	"siyahsensei/wallet-service/infrastructure/persistence/assetrepo"
	customLogger "siyahsensei/wallet-service/internal/pkg/logger"
	"siyahsensei/wallet-service/internal/platform/database"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize logger
	customLogger.InitLogger(customLogger.Config{
		LogLevel: "debug",
		Pretty:   config.Environment == "development",
	})

	customLogger.Info("Starting wallet worker...", map[string]interface{}{
		"environment": config.Environment,
	})

	// Connect to database
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

	// Initialize repositories
	accountRepo := accountrepo.NewPostgresRepository(db)
	assetRepo := assetrepo.NewPostgresRepository(db)

	// Initialize services
	accountService := account.NewService(accountRepo)
	assetService := asset.NewService(assetRepo)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the workers
	go runAssetPriceUpdater(ctx, assetService)
	go runAccountSynchronizer(ctx, accountService)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	customLogger.Info("Worker shutting down...")
	cancel()
	time.Sleep(time.Second)
	customLogger.Info("Worker stopped")
}

func runAssetPriceUpdater(ctx context.Context, assetService *asset.Service) {
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	customLogger.Info("Asset price updater started")

	for {
		select {
		case <-ctx.Done():
			customLogger.Info("Asset price updater stopped")
			return
		case <-ticker.C:
			updateAssetPrices(ctx, assetService)
		}
	}
}

func updateAssetPrices(ctx context.Context, assetService *asset.Service) {
	customLogger.Info("Updating asset prices...")
	// In a real application, you would:
	// 1. Fetch all assets that need price updates
	// 2. Query external APIs to get current prices
	// 3. Update asset prices in the database
	customLogger.Info("Asset price update completed")
}

func runAccountSynchronizer(ctx context.Context, accountService *account.Service) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	customLogger.Info("Account synchronizer started")

	for {
		select {
		case <-ctx.Done():
			customLogger.Info("Account synchronizer stopped")
			return
		case <-ticker.C:
			synchronizeAccounts(ctx, accountService)
		}
	}
}

func synchronizeAccounts(ctx context.Context, accountService *account.Service) {
	customLogger.Info("Synchronizing accounts...")
	// In a real application, you would:
	// 1. Fetch all accounts with API connections
	// 2. For each account, call the appropriate external API
	// 3. Update account balances and related data
	customLogger.Info("Account synchronization completed")
}
