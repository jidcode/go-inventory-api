package categories

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/utils"
)

type CategoryController struct {
	repo *CategoryRepository
}

func NewCategoryController(categoryRepo *CategoryRepository) *CategoryController {
	return &CategoryController{repo: categoryRepo}
}

func (ctrl *CategoryController) ListCategories(ctx echo.Context) error {
	inventoryId, err := uuid.Parse(ctx.Param("inventoryId"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid inventory ID")
	}

	categories, err := ctrl.repo.ListCategories(inventoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch categories",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, categories)
}

func (ctrl *CategoryController) GetCategory(ctx echo.Context) error {
	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	category, err := ctrl.repo.GetCategory(categoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, category)
}

func (ctrl *CategoryController) CreateCategory(ctx echo.Context) error {
	var input domain.CategoryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	newCategory := input.ToCreateCategoryRequest()

	err := ctrl.repo.CreateCategory(newCategory)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, newCategory)
}

func (ctrl *CategoryController) EditCategory(ctx echo.Context) error {
	var input domain.CategoryRequest
	input.Sanitize()

	if err := utils.BindAndValidateInput(ctx, &input); err != nil {
		return err
	}

	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	existingCategory, err := ctrl.repo.GetCategory(categoryId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Category not found")
	}

	updatedCategory := input.ToEditCategoryRequest(existingCategory)

	err = ctrl.repo.EditCategory(updatedCategory)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update category",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, updatedCategory)
}

func (ctrl *CategoryController) DeleteCategory(ctx echo.Context) error {
	categoryId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid category ID")
	}

	err = ctrl.repo.DeleteCategory(categoryId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete category",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
