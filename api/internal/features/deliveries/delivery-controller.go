package deliveries

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

type DeliveryController struct {
	repo *DeliveryRepository
}

func NewDeliveryController(deliveryRepo *DeliveryRepository) *DeliveryController {
	return &DeliveryController{repo: deliveryRepo}
}

func (ctrl *DeliveryController) ListDeliveries(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	deliveries, err := ctrl.repo.ListDeliveries(inventoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch deliveries",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, deliveries)
}

func (ctrl *DeliveryController) GetDelivery(ctx echo.Context) error {
	deliveryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid delivery ID")
	}

	delivery, err := ctrl.repo.GetDelivery(deliveryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve delivery",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, delivery)
}

func (ctrl *DeliveryController) CreateDelivery(ctx echo.Context) error {
	var input domain.DeliveryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newDelivery := input.ToCreateDeliveryRequest()

	if err := ctrl.repo.CreateDelivery(newDelivery); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create delivery",
			"details": err.Error(),
		})
	}

	delivery, err := ctrl.repo.GetDelivery(newDelivery.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve created delivery",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, delivery)
}

func (ctrl *DeliveryController) UpdateDelivery(ctx echo.Context) error {
	var input domain.DeliveryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	deliveryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid delivery ID")
	}

	existingDelivery, err := ctrl.repo.GetDelivery(deliveryId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Delivery not found")
	}

	updatedDelivery := input.ToUpdateDeliveryRequest(existingDelivery)

	if err := ctrl.repo.UpdateDelivery(updatedDelivery); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update delivery",
			"details": err.Error(),
		})
	}

	delivery, err := ctrl.repo.GetDelivery(updatedDelivery.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve updated delivery",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, delivery)
}

func (ctrl *DeliveryController) DeleteDelivery(ctx echo.Context) error {
	deliveryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid delivery ID")
	}

	if err := ctrl.repo.DeleteDelivery(deliveryId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete delivery",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
