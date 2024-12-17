package storages

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type StorageHandler struct {
	repo *StorageRepository
}

func NewStorageHandler(storageRepo *StorageRepository) *StorageHandler {
	return &StorageHandler{repo: storageRepo}
}

func (handler *StorageHandler) ListStorages(ctx echo.Context) error {
	inventoryID, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	storages, err := handler.repo.ListStorages(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch storages",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, storages)
}

func (handler *StorageHandler) GetStorage(ctx echo.Context) error {
	storageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	storage, err := handler.repo.GetStorage(storageID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, storage)
}

func (handler *StorageHandler) CreateStorage(ctx echo.Context) error {
	var input models.StorageInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage input")
	}

	input.Sanitize()

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	newStorage := input.ToStorage()

	if err := handler.repo.CreateStorage(&newStorage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, newStorage)
}

func (handler *StorageHandler) UpdateStorage(ctx echo.Context) error {
	storageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	var input models.StorageInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	input.Sanitize()

	existingStorage, err := handler.repo.GetStorage(storageID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Storage not found")
	}

	existingStorage.Name = input.Name
	existingStorage.Location = input.Location
	existingStorage.Capacity = input.Capacity
	existingStorage.InventoryID = input.InventoryID
	existingStorage.UpdatedAt = time.Now()

	if err := handler.repo.UpdateStorage(&existingStorage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingStorage)
}

func (handler *StorageHandler) DeleteStorage(ctx echo.Context) error {
	storageID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	if err := handler.repo.DeleteStorage(storageID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete storage",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
