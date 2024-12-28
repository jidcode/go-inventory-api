package products

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

type ProductController struct {
	repo *ProductRepository
}

func NewProductController(productRepo *ProductRepository) *ProductController {
	return &ProductController{repo: productRepo}
}

func (ctrl *ProductController) ListProducts(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	products, err := ctrl.repo.ListProducts(inventoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, products)
}

func (ctrl *ProductController) GetProduct(ctx echo.Context) error {
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product, err := ctrl.repo.GetProductWithRelations(productId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, product)
}

func (ctrl *ProductController) CreateProduct(ctx echo.Context) error {
	var input domain.ProductRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newProduct := input.ToCreateProductRequest()

	err := ctrl.repo.CreateProduct(newProduct, input.Categories, input.Storages, input.Images)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
	}

	product, err := ctrl.repo.GetProductWithRelations(newProduct.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, product)
}

func (ctrl *ProductController) EditProduct(ctx echo.Context) error {
	var input domain.ProductRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	existingProduct, err := ctrl.repo.GetProduct(productId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Product not found")
	}

	updatedProduct := input.ToEditProductRequest(existingProduct)

	err = ctrl.repo.EditProduct(updatedProduct, input.Categories, input.Storages, input.Images)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update product",
			"details": err.Error(),
		})
	}

	product, err := ctrl.repo.GetProductWithRelations(updatedProduct.Id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve updated product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, product)
}

func (ctrl *ProductController) DeleteProduct(ctx echo.Context) error {
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	err = ctrl.repo.DeleteProduct(productId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete product",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
