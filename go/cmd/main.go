package main

import (
	"log"

	"github.com/ventry/cmd/server"
	"github.com/ventry/database"
	"github.com/ventry/internal/domain/config"
	"github.com/ventry/internal/features/auth"
	"github.com/ventry/internal/features/categories"
	"github.com/ventry/internal/features/inventories"
	"github.com/ventry/internal/features/products"
	"github.com/ventry/internal/features/storages"
)

func main() {
	config := config.LoadEnv()

	// Initialize database connections
	db := database.Connect(*config)
	defer db.Close()

	// Database Repositories
	authRepo := auth.NewAuthRepository(db)
	inventoryRepo := inventories.NewInventoryRepository(db)
	productRepo := products.NewProductRepository(db)
	imagesRepo := products.NewImagesRepository(db)
	categoryRepo := categories.NewCategoryRepository(db)
	storageRepo := storages.NewStorageRepository(db)
	unitsRepo := storages.NewUnitsRepository(db)
	// orderRepo := orders.NewOrderRepository(db)
	// orderItemRepo := orders.NewOrderItemRepository(db)

	// Declare dependencies
	dependencies := server.ServerDependencies{
		AuthService:      auth.NewAuthService(authRepo, config),
		AuthHandler:      auth.NewAuthHandler(auth.NewAuthService(authRepo, config)),
		InventoryHandler: inventories.NewInventoryHandler(inventoryRepo),
		ProductHandler:   products.NewProductHandler(productRepo),
		ImagesHandler:    products.NewImagesHandler(imagesRepo),
		CategoryHandler:  categories.NewCategoryHandler(categoryRepo),
		StorageHandler:   storages.NewStorageHandler(storageRepo),
		UnitsHandler:     storages.NewUnitsHandler(unitsRepo),
		// OrderHandler:     orders.NewOrderHandler(orderRepo),
		// OrderItemHandler: orders.NewOrderItemHandler(orderItemRepo),
	}

	e := server.Run(dependencies)

	log.Println("Server running on port: 5000")
	log.Fatal(e.Start(":5000"))
}
