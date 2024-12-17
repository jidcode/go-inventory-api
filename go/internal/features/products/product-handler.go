package products

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type ProductHandler struct {
	repo *ProductRepository
}

func NewProductHandler(repo *ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (handler *ProductHandler) ListProducts(ctx echo.Context) error {
	inventoryID := ctx.Param("inventoryId")
	id, err := uuid.Parse(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	products, err := handler.repo.ListProducts(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, products)
}

func (handler *ProductHandler) GetProduct(ctx echo.Context) error {
	productID := ctx.Param("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	product, err := handler.repo.GetProductWithRelations(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, product)
}

func (handler *ProductHandler) CreateProduct(ctx echo.Context) error {
	var input models.ProductInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product input")
	}

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Validation failed")
	}

	input.Sanitize()

	// Initialize the new product object
	newProduct := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Code:        input.Code,
		Cost:        input.Cost,
		Price:       input.Price,
		Quantity:    input.Quantity,
		InventoryID: input.InventoryID,
	}

	// Create the product in the database
	if err := handler.repo.CreateProduct(&newProduct, input.Categories, input.Storages, input.Images); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
	}

	// Fetch the newly created product with categories, storages, and images
	productWithDetails, err := handler.repo.GetProductWithRelations(newProduct.ID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, productWithDetails)
}

func (handler *ProductHandler) UpdateProduct(ctx echo.Context) error {
	productID := ctx.Param("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	var input models.ProductInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	input.Sanitize()

	existingProduct, err := handler.repo.GetProduct(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Product not found")
	}

	existingProduct.Name = input.Name
	existingProduct.Description = input.Description
	existingProduct.Code = input.Code
	existingProduct.Cost = input.Cost
	existingProduct.Price = input.Price
	existingProduct.Quantity = input.Quantity
	existingProduct.InventoryID = input.InventoryID

	if err := handler.repo.UpdateProduct(&existingProduct); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update product",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingProduct)
}

func (handler *ProductHandler) DeleteProduct(ctx echo.Context) error {
	productID := ctx.Param("id")
	id, err := uuid.Parse(productID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	if err := handler.repo.DeleteProduct(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete product",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

