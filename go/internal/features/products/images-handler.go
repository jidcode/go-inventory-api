package products

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ventry/internal/domain/models"
)

type ImagesHandler struct {
	repo *ImagesRepository
}

func NewImagesHandler(imageRepo *ImagesRepository) *ImagesHandler {
	return &ImagesHandler{repo: imageRepo}
}

func (handler *ImagesHandler) ListImages(ctx echo.Context) error {
	productID := ctx.Param("productId")
	id, err := uuid.Parse(productID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid product ID")
	}

	images, err := handler.repo.ListImages(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to fetch images",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, images)
}

func (handler *ImagesHandler) GetImage(ctx echo.Context) error {
	imageID := ctx.Param("id")
	id, err := uuid.Parse(imageID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid image ID")
	}

	image, err := handler.repo.GetImage(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to retrieve image",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, image)
}

func (handler *ImagesHandler) CreateImage(ctx echo.Context) error {
	var input models.ImageInput
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid image input")
	}

	newImage := models.Image{
		URL:       input.URL,
		ProductID: input.ProductID,
		IsPrimary: input.IsPrimary,
	}

	if err := handler.repo.CreateImage(&newImage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to create image",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, newImage)
}

func (handler *ImagesHandler) UpdateImage(ctx echo.Context) error {
	imageID := ctx.Param("id")
	id, err := uuid.Parse(imageID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid image ID")
	}

	var input models.ImageInput
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid input")
	}

	existingImage, err := handler.repo.GetImage(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, "Image not found")
	}

	existingImage.URL = input.URL
	existingImage.ProductID = input.ProductID
	existingImage.IsPrimary = input.IsPrimary

	if err := handler.repo.UpdateImage(&existingImage); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to update image",
			"details": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, existingImage)
}

func (handler *ImagesHandler) DeleteImage(ctx echo.Context) error {
	imageID := ctx.Param("id")
	id, err := uuid.Parse(imageID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "Invalid image ID")
	}

	if err := handler.repo.DeleteImage(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error":   "Failed to delete image",
			"details": err.Error(),
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}
