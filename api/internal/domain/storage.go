// domain/storage.go
package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	Id          uuid.UUID     `db:"id" json:"id"`
	Name        string        `db:"name" json:"name"`
	Location    string        `db:"location" json:"location"`
	Capacity    int           `db:"capacity" json:"capacity"`
	InventoryId uuid.UUID     `db:"inventory_id" json:"inventoryId"`
	CreatedAt   time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updatedAt"`
	Units       []StorageUnit `json:"units,omitempty"`
}

// DTOs
type StorageRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=100"`
	Location    string    `json:"location" validate:"required"`
	Capacity    int       `json:"capacity" validate:"required,min=1"`
	InventoryId uuid.UUID `json:"inventoryId" validate:"required"`
	UnitCount   int       `json:"unit_count"`
}

func (req *StorageRequest) ToCreateStorageRequest() *Storage {
	return &Storage{
		Id:          uuid.New(),
		Name:        req.Name,
		Location:    req.Location,
		Capacity:    req.Capacity,
		InventoryId: req.InventoryId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (req *StorageRequest) ToUpdateStorageRequest(existing *Storage) *Storage {
	existing.Name = req.Name
	existing.Location = req.Location
	existing.Capacity = req.Capacity
	existing.InventoryId = req.InventoryId
	existing.UpdatedAt = time.Now()
	return existing
}

func (req *StorageRequest) Sanitize() {
	req.Name = strings.TrimSpace(req.Name)
	req.Location = strings.TrimSpace(req.Location)
}
