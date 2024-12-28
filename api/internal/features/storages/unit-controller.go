package storages

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

func (ctrl *StorageController) ListUnits(ctx echo.Context) error {
	storageId, err := uuid.Parse(ctx.Param("storageId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	units, err := ctrl.repo.ListUnits(storageId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch units",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, units)
}

func (ctrl *StorageController) GetUnit(ctx echo.Context) error {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	unit, err := ctrl.repo.GetUnit(unitId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, unit)
}

func (ctrl *StorageController) CreateUnit(ctx echo.Context) error {
	var input domain.StorageUnitRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newUnit := input.ToCreateStorageUnitRequest()

	err := ctrl.repo.CreateUnit(newUnit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create unit",
			"details": err.Error(),
		})
	}

	unit, err := ctrl.repo.GetUnit(newUnit.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve created unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, unit)
}

func (ctrl *StorageController) EditUnit(ctx echo.Context) error {
	var input domain.StorageUnitRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	existingUnit, err := ctrl.repo.GetUnit(unitId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Unit not found")
	}

	updatedUnit := input.ToUpdateStorageUnitRequest(existingUnit)

	err = ctrl.repo.EditUnit(updatedUnit)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update unit",
			"details": err.Error(),
		})
	}

	unit, err := ctrl.repo.GetUnit(updatedUnit.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve updated unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, unit)
}

func (ctrl *StorageController) DeleteUnit(ctx echo.Context) error {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	err = ctrl.repo.DeleteUnit(unitId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete unit",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// HELPERS
func (ctrl *StorageController) AddItemToUnit(ctx echo.Context) error {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	var input domain.UnitItemRequest
	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	err = ctrl.repo.AddItemToUnit(unitId, input.ProductId, input.Quantity)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to add item to unit",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (ctrl *StorageController) RemoveItemFromUnit(ctx echo.Context) error {
	unitId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	var input domain.UnitItemRequest
	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	err = ctrl.repo.RemoveItemFromUnit(unitId, input.Quantity)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to remove item from unit",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
