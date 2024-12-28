package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ventry/internal/features/categories"
	"github.com/ventry/internal/features/deliveries"
	"github.com/ventry/internal/features/inventories"
	"github.com/ventry/internal/features/products"
	"github.com/ventry/internal/features/sales"
	"github.com/ventry/internal/features/storages"
	"github.com/ventry/internal/pkg/auth"
	"github.com/ventry/internal/pkg/logger"
	"github.com/ventry/internal/router"
	cv "github.com/ventry/internal/utils"
)

type ServerDependencies struct {
	AuthService         *auth.AuthService
	AuthHandler         *auth.AuthHandler
	InventoryController *inventories.InventoryController
	StorageController   *storages.StorageController
	ProductController   *products.ProductController
	CategoryController  *categories.CategoryController
	DeliveryController  *deliveries.DeliveryController
	SaleController      *sales.SaleController
}

func Run(deps ServerDependencies) *echo.Echo {
	e := echo.New()

	e.Validator = &cv.CustomValidator{Validator: validator.New()}

	e.Use(logger.RequestMiddleware)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
	}))

	e.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]string{"status": "OK"})
	})

	//Routing
	router.AuthRoutes(e, deps.AuthHandler, deps.AuthService)
	router.InventoryRoutes(e, *deps.InventoryController, *deps.AuthService)
	router.StorageRoutes(e, *deps.StorageController, *deps.AuthService)
	router.ProductRoutes(e, *deps.ProductController, *deps.AuthService)
	router.CategoryRoutes(e, *deps.CategoryController, *deps.AuthService)
	router.DeliveryRoutes(e, *deps.DeliveryController, *deps.AuthService)
	router.SaleRoutes(e, *deps.SaleController, *deps.AuthService)

	return e
}
