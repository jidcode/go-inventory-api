package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type DeliveryStatus string

const (
	DeliveryStatusPending    DeliveryStatus = "pending"
	DeliveryStatusProcessing DeliveryStatus = "processing"
	DeliveryStatusShipped    DeliveryStatus = "shipped"
	DeliveryStatusDelivered  DeliveryStatus = "delivered"
	DeliveryStatusCancelled  DeliveryStatus = "cancelled"
)

type Delivery struct {
	Id               uuid.UUID      `db:"id" json:"id"`
	InventoryId      uuid.UUID      `db:"inventory_id" json:"inventoryId"`
	Status           DeliveryStatus `db:"status" json:"status"`
	OrderDate        time.Time      `db:"order_date" json:"orderDate"`
	DeliveredAt      *time.Time     `db:"delivered_at" json:"deliveredAt"`
	RecipientName    string         `db:"recipient_name" json:"recipientName"`
	RecipientAddress string         `db:"recipient_address" json:"recipientAddress"`
	RecipientPhone   string         `db:"recipient_phone" json:"recipientPhone"`
	TrackingNumber   string         `db:"tracking_number" json:"trackingNumber"`
	Note             string         `db:"note" json:"note"`
	CreatedAt        time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time      `db:"updated_at" json:"updatedAt"`
	Items            []DeliveryItem `json:"items"`
}

type DeliveryItem struct {
	Id         uuid.UUID `db:"id" json:"id"`
	DeliveryId uuid.UUID `db:"delivery_id" json:"deliveryId"`
	ProductId  uuid.UUID `db:"product_id" json:"productId"`
	Quantity   int       `db:"quantity" json:"quantity"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	Product    *Product  `json:"product,omitempty"`
}

// DTOs
type DeliveryRequest struct {
	InventoryId      uuid.UUID             `json:"inventoryId" validate:"required"`
	Status           DeliveryStatus        `json:"status" validate:"required,oneof=pending processing shipped delivered cancelled"`
	OrderDate        time.Time             `json:"orderDate" validate:"required"`
	DeliveredAt      *time.Time            `json:"deliveredAt"`
	RecipientName    string                `json:"recipientName" validate:"required"`
	RecipientAddress string                `json:"recipientAddress" validate:"required"`
	RecipientPhone   string                `json:"recipientPhone" validate:"required"`
	TrackingNumber   string                `json:"trackingNumber"`
	Note             string                `json:"note"`
	Items            []DeliveryItemRequest `json:"items" validate:"required,min=1"`
}

type DeliveryItemRequest struct {
	ProductId uuid.UUID `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

func (req *DeliveryRequest) ToCreateDeliveryRequest() *Delivery {
	delivery := &Delivery{
		Id:               uuid.New(),
		InventoryId:      req.InventoryId,
		Status:           req.Status,
		OrderDate:        req.OrderDate,
		DeliveredAt:      req.DeliveredAt,
		RecipientName:    req.RecipientName,
		RecipientAddress: req.RecipientAddress,
		RecipientPhone:   req.RecipientPhone,
		TrackingNumber:   req.TrackingNumber,
		Note:             req.Note,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	delivery.Items = make([]DeliveryItem, len(req.Items))
	for i, item := range req.Items {
		delivery.Items[i] = DeliveryItem{
			Id:         uuid.New(),
			DeliveryId: delivery.Id,
			ProductId:  item.ProductId,
			Quantity:   item.Quantity,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	return delivery
}

func (req *DeliveryRequest) ToUpdateDeliveryRequest(existingDelivery *Delivery) *Delivery {
	existingDelivery.Status = req.Status
	existingDelivery.OrderDate = req.OrderDate
	existingDelivery.DeliveredAt = req.DeliveredAt
	existingDelivery.RecipientName = req.RecipientName
	existingDelivery.RecipientAddress = req.RecipientAddress
	existingDelivery.RecipientPhone = req.RecipientPhone
	existingDelivery.TrackingNumber = req.TrackingNumber
	existingDelivery.Note = req.Note
	existingDelivery.UpdatedAt = time.Now()

	existingDelivery.Items = make([]DeliveryItem, len(req.Items))
	for i, item := range req.Items {
		existingDelivery.Items[i] = DeliveryItem{
			Id:         uuid.New(),
			DeliveryId: existingDelivery.Id,
			ProductId:  item.ProductId,
			Quantity:   item.Quantity,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
	}

	return existingDelivery
}

func (req *DeliveryRequest) Sanitize() {
	req.RecipientName = strings.TrimSpace(req.RecipientName)
	req.RecipientAddress = strings.TrimSpace(req.RecipientAddress)
	req.RecipientPhone = strings.TrimSpace(req.RecipientPhone)
	req.TrackingNumber = strings.TrimSpace(req.TrackingNumber)
	req.Note = strings.TrimSpace(req.Note)
}
