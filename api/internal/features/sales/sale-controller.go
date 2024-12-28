package sales

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

type SaleController struct {
	repo *SaleRepository
}

func NewSaleController(saleRepo *SaleRepository) *SaleController {
	return &SaleController{repo: saleRepo}
}

func (ctrl *SaleController) ListSales(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	sales, err := ctrl.repo.ListSales(inventoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch sales",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, sales)
}

func (ctrl *SaleController) GetSale(ctx echo.Context) error {
	saleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid sale ID")
	}

	sale, err := ctrl.repo.GetSale(saleId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve sale",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, sale)
}

func (ctrl *SaleController) CreateSale(ctx echo.Context) error {
	var input domain.SaleCreateRequest
	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	sale := input.ToSale()

	err := ctrl.repo.CreateSale(sale)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create sale",
			"details": err.Error(),
		})
	}

	// Fetch the complete sale with items
	createdSale, err := ctrl.repo.GetSale(sale.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve created sale",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, createdSale)
}

func (ctrl *SaleController) EditSale(ctx echo.Context) error {
	var input domain.SaleRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	saleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid sale ID")
	}

	existingSale, err := ctrl.repo.GetSale(saleId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Sale not found")
	}

	updatedSale := input.ToUpdateSaleRequest(existingSale)

	err = ctrl.repo.EditSale(updatedSale)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update sale",
			"details": err.Error(),
		})
	}

	sale, err := ctrl.repo.GetSale(updatedSale.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve updated sale",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, sale)
}

func (ctrl *SaleController) DeleteSale(ctx echo.Context) error {
	saleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid sale ID")
	}

	err = ctrl.repo.DeleteSale(saleId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete sale",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// HELPERS
func (ctrl *SaleController) AddItemToSale(ctx echo.Context) error {
	saleId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid sale ID")
	}

	var input domain.SaleItemRequest
	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	item := &domain.SaleItem{
		Id:        uuid.New(),
		SaleId:    saleId,
		ProductId: input.ProductId,
		Quantity:  input.Quantity,
		UnitPrice: input.UnitPrice,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ctrl.repo.AddItemToSale(item)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to add item to sale",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

func (ctrl *SaleController) RemoveItemFromSale(ctx echo.Context) error {
	saleId, err := uuid.Parse(ctx.Param("saleId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid sale ID")
	}

	itemId, err := uuid.Parse(ctx.Param("itemId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid item ID")
	}

	err = ctrl.repo.RemoveItemFromSale(saleId, itemId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to remove item from sale",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
