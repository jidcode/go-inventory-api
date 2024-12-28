package domain

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	Id        uuid.UUID `db:"id" json:"id"`
	URL       string    `db:"url" json:"url"`
	ProductId uuid.UUID `db:"product_id" json:"productId"`
	IsPrimary bool      `db:"is_primary" json:"isPrimary"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type ImageRequest struct {
	URL       string    `json:"url" validate:"required,url"`
	ProductId uuid.UUID `json:"productId" validate:"required"`
	IsPrimary bool      `json:"isPrimary,omitempty"`
}
