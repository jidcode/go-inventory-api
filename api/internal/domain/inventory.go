package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	UserId      uuid.UUID `db:"user_id" json:"userId"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// DTOs
type InventoryRequest struct {
	Name        string    `json:"name" validate:"required,min=3,max=50"`
	Description string    `json:"description"`
	UserId      uuid.UUID `json:"userId" validate:"required"`
}

type InventoryResponse struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Name        string    `json:"name" validate:"required,min=3,max=50"`
	Description string    `json:"description"`
	UserId      uuid.UUID `json:"userId" validate:"required"`
}

func (req *InventoryRequest) ToCreateInventoryRequest() *Inventory {
	return &Inventory{
		Id:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		UserId:      req.UserId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (req *InventoryRequest) ToUpdateInventoryRequest(existingInventory *Inventory) *Inventory {
	existingInventory.Name = req.Name
	existingInventory.Description = req.Description
	existingInventory.UserId = req.UserId
	existingInventory.UpdatedAt = time.Now()

	return existingInventory
}

func (req *InventoryRequest) Sanitize() {
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
}
