package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ventry/internal/features/auth"
	"github.com/ventry/internal/features/categories"
	"github.com/ventry/internal/features/inventories"
	"github.com/ventry/internal/features/products"
	"github.com/ventry/internal/features/storages"
	cv "github.com/ventry/internal/utils"
)

type ServerDependencies struct {
	AuthService      *auth.AuthService
	AuthHandler      *auth.AuthHandler
	InventoryHandler *inventories.InventoryHandler
	ProductHandler   *products.ProductHandler
	ImagesHandler    *products.ImagesHandler
	CategoryHandler  *categories.CategoryHandler
	StorageHandler   *storages.StorageHandler
	UnitsHandler     *storages.UnitsHandler
	// OrderHandler     *orders.OrderHandler
	// OrderItemHandler *orders.OrderItemHandler
}

func Run(deps ServerDependencies) *echo.Echo {
	e := echo.New()

	e.Validator = &cv.CustomValidator{Validator: validator.New()}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE},
	}))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]string{"status": "OK"})
	})

	//Routing
	auth.AuthRouter(e, deps.AuthHandler, deps.AuthService)
	inventories.InventoryRouter(e, *deps.InventoryHandler, *deps.AuthService)
	products.ProductRouter(e, deps.ProductHandler, deps.ImagesHandler, *deps.AuthService)
	categories.CategoryRouter(e, deps.CategoryHandler, *deps.AuthService)
	storages.StorageRouter(e, deps.StorageHandler, deps.UnitsHandler, *deps.AuthService)
	// orders.OrderRouter(e, deps.OrderHandler, deps.OrderItemHandler, *deps.AuthService)

	return e
}
