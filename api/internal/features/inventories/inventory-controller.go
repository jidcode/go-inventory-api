package inventories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/pkg/errors"
	"github.com/ventry/internal/pkg/logger"
	"github.com/ventry/internal/utils"
)

type InventoryController struct {
	repo *InventoryRepository
}

func NewInventoryController(inventoryRepo *InventoryRepository) *InventoryController {
	return &InventoryController{repo: inventoryRepo}
}

func (control *InventoryController) ListInventories(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*domain.User)
	if !ok {
		return errors.New(errors.Unauthorized, "User not authenticated", http.StatusUnauthorized)
	}

	inventories, err := control.repo.ListInventories(user.Id)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Failed to fetch inventories",
			logger.Field{Key: "user_id", Value: user.Id})
		return errors.Send(ctx, err)
	}

	logger.Info(ctx.Request().Context(), "Successfully fetched inventories",
		logger.Field{Key: "user_id", Value: user.Id},
		logger.Field{Key: "count", Value: len(inventories)})

	return ctx.JSON(http.StatusOK, inventories)
}

func (control *InventoryController) GetInventory(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errors.Send(ctx, err)
		return errors.ValidationError("Invalid inventory ID")
	}

	inventory, err := control.repo.GetInventory(inventoryId)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Failed to retrieve inventory",
			logger.Field{Key: "inventory_id", Value: inventoryId})
		return errors.Send(ctx, err)
	}

	logger.Info(ctx.Request().Context(), "Successfully retrieved inventory",
		logger.Field{Key: "inventory_id", Value: inventoryId})

	return ctx.JSON(http.StatusOK, inventory)
}

func (control *InventoryController) CreateInventory(ctx echo.Context) error {
	var input domain.InventoryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newInventory := input.ToCreateInventoryRequest()

	err := control.repo.CreateInventory(newInventory)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Failed to create inventory",
			logger.Field{Key: "inventory_name", Value: newInventory.Name})
		return errors.Send(ctx, err)
	}

	logger.Info(ctx.Request().Context(), "Successfully created inventory",
		logger.Field{Key: "inventory_id", Value: newInventory.Id},
		logger.Field{Key: "inventory_name", Value: newInventory.Name})

	return ctx.JSON(http.StatusCreated, newInventory)
}

func (control *InventoryController) EditInventory(ctx echo.Context) error {
	var input domain.InventoryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	inventoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return errors.ValidationError("Invalid inventory ID")
	}

	existingInventory, err := control.repo.GetInventory(inventoryId)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Inventory not found",
			logger.Field{Key: "inventory_id", Value: inventoryId})
		return errors.Send(ctx, err)
	}

	updatedInventory := input.ToUpdateInventoryRequest(&existingInventory)

	err = control.repo.EditInventory(updatedInventory)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Failed to update inventory",
			logger.Field{Key: "inventory_id", Value: inventoryId})
		return errors.Send(ctx, err)
	}

	logger.Info(ctx.Request().Context(), "Successfully updated inventory",
		logger.Field{Key: "inventory_id", Value: inventoryId},
		logger.Field{Key: "inventory_name", Value: updatedInventory.Name})

	return ctx.JSON(http.StatusOK, updatedInventory)
}

func (control *InventoryController) DeleteInventory(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return errors.ValidationError("Invalid inventory ID")
	}

	err = control.repo.DeleteInventory(inventoryId)
	if err != nil {
		logger.Error(ctx.Request().Context(), err, "Failed to delete inventory",
			logger.Field{Key: "inventory_id", Value: inventoryId})
		return errors.Send(ctx, err)
	}

	logger.Info(ctx.Request().Context(), "Successfully deleted inventory",
		logger.Field{Key: "inventory_id", Value: inventoryId})

	return ctx.NoContent(http.StatusNoContent)
}
