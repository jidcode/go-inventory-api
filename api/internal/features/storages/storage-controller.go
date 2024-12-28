package storages

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

type StorageController struct {
	repo *StorageRepository
}

func NewStorageController(storageRepo *StorageRepository) *StorageController {
	return &StorageController{repo: storageRepo}
}

func (ctrl *StorageController) ListStorages(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	storages, err := ctrl.repo.ListStorages(inventoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch storages",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, storages)
}

func (ctrl *StorageController) GetStorage(ctx echo.Context) error {
	storageId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	storage, err := ctrl.repo.GetStorage(storageId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, storage)
}

func (ctrl *StorageController) CreateStorage(ctx echo.Context) error {
	var input domain.StorageRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newStorage := input.ToCreateStorageRequest()

	err := ctrl.repo.CreateStorage(newStorage)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create storage",
			"details": err.Error(),
		})
	}

	storage, err := ctrl.repo.GetStorage(newStorage.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve created storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, storage)
}

func (ctrl *StorageController) EditStorage(ctx echo.Context) error {
	var input domain.StorageRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	storageId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	existingStorage, err := ctrl.repo.GetStorage(storageId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Storage not found")
	}

	updatedStorage := input.ToUpdateStorageRequest(existingStorage)

	err = ctrl.repo.EditStorage(updatedStorage)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update storage",
			"details": err.Error(),
		})
	}

	storage, err := ctrl.repo.GetStorage(updatedStorage.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve updated storage",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, storage)
}

func (ctrl *StorageController) DeleteStorage(ctx echo.Context) error {
	storageId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	err = ctrl.repo.DeleteStorage(storageId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete storage",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
