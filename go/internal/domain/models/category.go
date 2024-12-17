package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	InventoryID uuid.UUID `db:"inventory_id" json:"inventoryID"`
	ProductID   uuid.UUID `db:"product_id" json:"productID"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type CategoryInput struct {
	Name        string    `json:"name" validate:"required,min=3,max=100"`
	InventoryID uuid.UUID `json:"inventoryID" validate:"required"`
	ProductID   uuid.UUID `json:"productID,omitempty"`
}

func (c *CategoryInput) Sanitize() {
	c.Name = strings.TrimSpace(c.Name)
}
