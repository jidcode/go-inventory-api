package main

import (
	"context"

	"github.com/ventry/cmd/server"
	"github.com/ventry/config"
	"github.com/ventry/database"
	"github.com/ventry/internal/features/categories"
	"github.com/ventry/internal/features/deliveries"
	"github.com/ventry/internal/features/inventories"
	"github.com/ventry/internal/features/products"
	"github.com/ventry/internal/features/sales"
	"github.com/ventry/internal/features/storages"
	"github.com/ventry/internal/pkg/auth"
	"github.com/ventry/internal/pkg/logger"
)

func main() {
	config := config.LoadEnv()

	// Initialize logger
	if err := logger.Init(config.Environment); err != nil {
		panic(err)
	}
	defer logger.Sync()

	// Initialize database connection
	db := database.Connect(*config)
	defer db.Close()

	// Add startup logging
	logger.Info(context.Background(), "Starting application...",
		logger.Field{Key: "environment", Value: config.Environment},
		logger.Field{Key: "database_url", Value: config.DatabaseUrl},
	)

	// Database Repositories
	authRepo := auth.NewAuthRepository(db)
	inventoryRepo := inventories.NewInventoryRepository(db)
	storageRepo := storages.NewStorageRepository(db)
	productRepo := products.NewProductRepository(db)
	categoryRepo := categories.NewCategoryRepository(db)
	deliveryRepo := deliveries.NewDeliveryRepository(db)
	saleRepo := sales.NewSaleRepository(db)

	// Declare dependencies
	dependencies := server.ServerDependencies{
		AuthService:         auth.NewAuthService(authRepo, config),
		AuthHandler:         auth.NewAuthHandler(auth.NewAuthService(authRepo, config)),
		InventoryController: inventories.NewInventoryController(inventoryRepo),
		StorageController:   storages.NewStorageController(storageRepo),
		ProductController:   products.NewProductController(productRepo),
		CategoryController:  categories.NewCategoryController(categoryRepo),
		DeliveryController:  deliveries.NewDeliveryController(deliveryRepo),
		SaleController:      sales.NewSaleController(saleRepo),
	}

	e := server.Run(dependencies)

	port := ":5000"
	logger.Info(context.Background(), "Server running on:",
		logger.Field{Key: "port", Value: port},
	)

	if err := e.Start(port); err != nil {
		logger.Fatal(context.Background(), "Server failed to start",
			logger.Field{Key: "error", Value: err.Error()},
		)
	}
}
