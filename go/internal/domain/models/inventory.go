package models

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	UserID      uuid.UUID `db:"user_id" json:"userID"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type InventoryInput struct {
	Name        string    `json:"name" validate:"required,min=3,max=50"`
	Description string    `json:"description"`
	UserID      uuid.UUID `json:"userID" validate:"required"`
}
