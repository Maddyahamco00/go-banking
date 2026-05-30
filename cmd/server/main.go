package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Provide your JWT as: `Authorization: Bearer <token>`

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-banking/docs"

	"go-banking/internal/config"
	"go-banking/internal/handler"
	"go-banking/internal/middleware"
	"go-banking/internal/pkg/cache"
	"go-banking/internal/pkg/database"
	"go-banking/internal/repository"
	"go-banking/internal/service"
)

func main() {
	// Load .env if it exists
	_ = godotenv.Load()

	// Load configuration
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize PostgreSQL
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL")

	// Initialize Redis
	redisClient, err := cache.NewRedisClient(&cfg.Redis)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v (continuing without cache)", err)
	} else {
		defer redisClient.Close()
		log.Println("Connected to Redis")
	}

	// Initialize repositories
	accountRepo := repository.NewAccountRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	idempotencyRepo := repository.NewIdempotencyRepository(db)
	walletRepo := repository.NewWalletRepository(db)

	// Initialize services
	ledgerSvc := service.NewLedgerService(db, accountRepo, ledgerRepo, transactionRepo)
	walletSvc := service.NewWalletService(walletRepo, accountRepo, ledgerSvc, transactionRepo)

	// Initialize handlers
	walletHandler := handler.NewWalletHandler(walletSvc)
	adminHandler := handler.NewAdminHandler(accountRepo, transactionRepo, walletRepo)

	// Initialize middleware
	idempotencyMiddleware := middleware.NewIdempotencyMiddleware(redisClient, idempotencyRepo, 24)
	rateLimiter := middleware.NewRateLimiter(&cfg.RateLimit)

	// Setup Gin router
	router := gin.Default()

	// Global middleware
	router.Use(rateLimiter.Handle())

	// Health check (no auth required)
	router.GET("/health", adminHandler.HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Wallet routes (JWT protected)
		wallet := v1.Group("/wallet")
		wallet.Use(middleware.JWTAuth(&cfg.JWT))
		wallet.Use(idempotencyMiddleware.Handle())
		{
			wallet.POST("/create", walletHandler.CreateWallet)
			wallet.GET("/:owner_id", walletHandler.GetWallet)
			wallet.GET("/:owner_id/transactions", walletHandler.GetWalletTransactions)
			wallet.GET("/:owner_id/ledger", walletHandler.GetWalletLedgerEntries)
			wallet.POST("/fund", walletHandler.FundWallet)
		}

		// Admin routes (JWT + admin role required)
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(&cfg.JWT))
		admin.Use(middleware.RequireRole("admin"))
		{
			admin.GET("/transactions", adminHandler.GetTransactions)
			admin.GET("/wallets", adminHandler.GetWallets)
			admin.GET("/loans", adminHandler.GetLoans)
		}
	}

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}
