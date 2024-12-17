package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	ID            uuid.UUID `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Location      string    `db:"location" json:"location"`
	Capacity      int       `db:"capacity" json:"capacity"`
	OccupiedSpace int       `db:"occupied_space" json:"occupiedSpace"`
	InventoryID   uuid.UUID `db:"inventory_id" json:"inventoryID"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
	Units         []Unit    `db:"units" json:"units"`
}

type StorageInput struct {
	Name          string    `json:"name" validate:"required,min=3,max=100"`
	Location      string    `json:"location"`
	Capacity      int       `json:"capacity" validate:"min=0"`
	InventoryID   uuid.UUID `json:"inventoryID" validate:"required"`
	OccupiedSpace int       `json:"occupiedSpace" validate:"min=0"`
}

func (s *StorageInput) Sanitize() {
	s.Name = strings.TrimSpace(s.Name)
	s.Location = strings.TrimSpace(s.Location)
}

func (s *StorageInput) ToStorage() Storage {
	return Storage{
		ID:            uuid.New(),
		Name:          s.Name,
		Location:      s.Location,
		Capacity:      s.Capacity,
		InventoryID:   s.InventoryID,
		OccupiedSpace: s.OccupiedSpace,
	}
}
