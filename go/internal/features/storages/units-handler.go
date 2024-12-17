package storages

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type UnitsHandler struct {
	repo *UnitsRepository
}

func NewUnitsHandler(unitRepo *UnitsRepository) *UnitsHandler {
	return &UnitsHandler{repo: unitRepo}
}

func (handler *UnitsHandler) ListUnits(ctx echo.Context) error {
	storageID := ctx.Param("storageId")

	id, err := uuid.Parse(storageID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid storage ID")
	}

	units, err := handler.repo.ListUnits(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch units",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, units)
}

func (handler *UnitsHandler) GetUnit(ctx echo.Context) error {
	unitID := ctx.Param("id")

	id, err := uuid.Parse(unitID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	unit, err := handler.repo.GetUnit(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, unit)
}

func (handler *UnitsHandler) CreateUnit(ctx echo.Context) error {
	var input models.UnitInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit input")
	}

	input.Sanitize()

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	newUnit := input.ToUnit()

	if err := handler.repo.CreateUnit(&newUnit); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, newUnit)
}

func (handler *UnitsHandler) UpdateUnit(ctx echo.Context) error {
	unitID := ctx.Param("id")

	id, err := uuid.Parse(unitID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	var input models.UnitInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	existingUnit, err := handler.repo.GetUnit(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Unit not found")
	}

	existingUnit.Label = input.Label
	existingUnit.StorageID = input.StorageID
	existingUnit.IsOccupied = input.IsOccupied
	existingUnit.Capacity = input.Capacity

	if err := handler.repo.UpdateUnit(&existingUnit); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update unit",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingUnit)
}

func (handler *UnitsHandler) DeleteUnit(ctx echo.Context) error {
	unitID := ctx.Param("id")

	id, err := uuid.Parse(unitID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid unit ID")
	}

	if err := handler.repo.DeleteUnit(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete unit",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
