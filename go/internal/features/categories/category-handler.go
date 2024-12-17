package categories

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type CategoryHandler struct {
	repo *CategoryRepository
}

func NewCategoryHandler(categoryRepo *CategoryRepository) *CategoryHandler {
	return &CategoryHandler{repo: categoryRepo}
}

func (handler *CategoryHandler) ListCategories(ctx echo.Context) error {
	inventoryID, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	categories, err := handler.repo.ListCategories(inventoryID)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch categories",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, categories)
}

func (handler *CategoryHandler) GetCategory(ctx echo.Context) error {
	categoryID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	category, err := handler.repo.GetCategory(categoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, "Category not found")
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, category)
}

func (handler *CategoryHandler) CreateCategory(ctx echo.Context) error {
	var input models.CategoryInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category input")
	}

	input.Sanitize()

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Category validation failed")
	}

	newCategory := models.Category{
		Name:        input.Name,
		InventoryID: input.InventoryID,
		ProductID:   input.ProductID,
	}

	if err := handler.repo.CreateCategory(&newCategory); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, newCategory)
}

func (handler *CategoryHandler) UpdateCategory(ctx echo.Context) error {
	categoryID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	var input models.CategoryInput

	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category input")
	}

	input.Sanitize()

	if err := ctx.Validate(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Category validation failed")
	}

	existingCategory, err := handler.repo.GetCategory(categoryID)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Category not found")
	}

	existingCategory.Name = input.Name
	existingCategory.InventoryID = input.InventoryID

	if err := handler.repo.UpdateCategory(&existingCategory); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingCategory)
}

func (handler *CategoryHandler) DeleteCategory(ctx echo.Context) error {
	categoryID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	if err := handler.repo.DeleteCategory(categoryID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete category",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
