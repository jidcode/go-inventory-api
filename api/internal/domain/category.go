package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventoryId"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// DTOs
type CategoryRequest struct {
	Name        string    `db:"name" json:"name"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventoryId"`
}

type CategoryResponse struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	InventoryId uuid.UUID `db:"inventory_id" json:"inventoryId"`
}

func (req *CategoryRequest) ToCreateCategoryRequest() *Category {
	return &Category{
		Id:          uuid.New(),
		Name:        req.Name,
		InventoryId: req.InventoryId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (req *CategoryRequest) ToEditCategoryRequest(existingCategory *Category) *Category {
	existingCategory.Name = req.Name
	existingCategory.InventoryId = req.InventoryId
	existingCategory.UpdatedAt = time.Now()

	return existingCategory
}

func (req *CategoryRequest) Sanitize() {
	req.Name = strings.TrimSpace(req.Name)
}
