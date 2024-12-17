package inventories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type InventoryHandler struct {
	repo *InventoryRepository
}

func NewInventoryHandler(inventoryRepo *InventoryRepository) *InventoryHandler {
	return &InventoryHandler{repo: inventoryRepo}
}

func (handler *InventoryHandler) ListInventories(ctx echo.Context) error {
	user, ok := ctx.Get("user").(*models.User)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "User not found")
	}

	inventories, err := handler.repo.ListInventories(user.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch inventories",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, inventories)
}

func (handler *InventoryHandler) GetInventory(ctx echo.Context) error {
	inventoryID := ctx.Param("id")

	id, err := uuid.Parse(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	inventory, err := handler.repo.GetInventory(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve inventory",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, inventory)
}

func (handler *InventoryHandler) CreateInventory(ctx echo.Context) error {
	var input models.InventoryInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory input")
	}

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	newInventory := models.Inventory{
		Name:        input.Name,
		Description: input.Description,
		UserID:      input.UserID,
	}

	if err := handler.repo.CreateInventory(&newInventory); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create inventory",
			"details": err.Error(),
		})

	}

	return ctx.JSON(http.StatusCreated, newInventory)
}

func (handler *InventoryHandler) UpdateInventory(ctx echo.Context) error {
	inventoryID := ctx.Param("id")

	id, err := uuid.Parse(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	var input models.InventoryInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	existingInventory, err := handler.repo.GetInventory(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Inventory not found")
	}

	existingInventory.Name = input.Name
	existingInventory.Description = input.Description
	existingInventory.UserID = input.UserID

	if err := handler.repo.UpdateInventory(&existingInventory); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update inventory",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingInventory)
}

func (handler *InventoryHandler) DeleteInventory(ctx echo.Context) error {
	inventoryID := ctx.Param("id")

	id, err := uuid.Parse(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	if err := handler.repo.DeleteInventory(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete inventory",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
