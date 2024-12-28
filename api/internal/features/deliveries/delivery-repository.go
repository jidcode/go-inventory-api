package deliveries

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

type DeliveryRepository struct {
	db *sqlx.DB
}

func NewDeliveryRepository(data *sqlx.DB) *DeliveryRepository {
	return &DeliveryRepository{db: data}
}

func (repo *DeliveryRepository) ListDeliveries(inventoryId uuid.UUID) ([]domain.Delivery, error) {
	deliveries := []domain.Delivery{}
	query := `SELECT * FROM deliveries WHERE inventory_id = $1 ORDER BY order_date DESC`

	if err := repo.db.Select(&deliveries, query, inventoryId); err != nil {
		return nil, err
	}

	// Fetch items for each delivery
	for i := range deliveries {
		items, err := repo.getDeliveryItems(deliveries[i].Id)
		if err != nil {
			return nil, err
		}
		deliveries[i].Items = items
	}

	return deliveries, nil
}

func (repo *DeliveryRepository) GetDelivery(deliveryId uuid.UUID) (*domain.Delivery, error) {
	var delivery domain.Delivery
	query := `SELECT * FROM deliveries WHERE id = $1`

	if err := repo.db.Get(&delivery, query, deliveryId); err != nil {
		return nil, err
	}

	items, err := repo.getDeliveryItems(deliveryId)
	if err != nil {
		return nil, err
	}
	delivery.Items = items

	return &delivery, nil
}

func (repo *DeliveryRepository) CreateDelivery(delivery *domain.Delivery) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert delivery
	deliveryQuery := `
		INSERT INTO deliveries (
			id, inventory_id, status, order_date, delivered_at,
			recipient_name, recipient_address, recipient_phone,
			tracking_number, note, created_at, updated_at
		) VALUES (
			:id, :inventory_id, :status, :order_date, :delivered_at,
			:recipient_name, :recipient_address, :recipient_phone,
			:tracking_number, :note, :created_at, :updated_at
		)
	`

	if _, err := tx.NamedExec(deliveryQuery, delivery); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert delivery items
	itemQuery := `
		INSERT INTO delivery_items (
			id, delivery_id, product_id, quantity, created_at, updated_at
		) VALUES (
			:id, :delivery_id, :product_id, :quantity, :created_at, :updated_at
		)
	`

	for _, item := range delivery.Items {
		if _, err := tx.NamedExec(itemQuery, item); err != nil {
			_ = tx.Rollback()
			return err
		}

		// Update product quantity
		updateQuery := `
			UPDATE products 
			SET quantity = quantity - $1, 
				updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
		`
		if _, err := tx.Exec(updateQuery, item.Quantity, item.ProductId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *DeliveryRepository) UpdateDelivery(delivery *domain.Delivery) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Get existing delivery items to restore quantities
	oldItems, err := repo.getDeliveryItems(delivery.Id)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Restore product quantities
	for _, item := range oldItems {
		restoreQuery := `
			UPDATE products 
			SET quantity = quantity + $1, 
				updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
		`
		if _, err := tx.Exec(restoreQuery, item.Quantity, item.ProductId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Update delivery
	deliveryQuery := `
		UPDATE deliveries SET
			status = :status,
			order_date = :order_date,
			delivered_at = :delivered_at,
			recipient_name = :recipient_name,
			recipient_address = :recipient_address,
			recipient_phone = :recipient_phone,
			tracking_number = :tracking_number,
			note = :note,
			updated_at = :updated_at
		WHERE id = :id
	`

	if _, err := tx.NamedExec(deliveryQuery, delivery); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Delete old items
	if _, err := tx.Exec(`DELETE FROM delivery_items WHERE delivery_id = $1`, delivery.Id); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert new items
	itemQuery := `
		INSERT INTO delivery_items (
			id, delivery_id, product_id, quantity, created_at, updated_at
		) VALUES (
			:id, :delivery_id, :product_id, :quantity, :created_at, :updated_at
		)
	`

	for _, item := range delivery.Items {
		if _, err := tx.NamedExec(itemQuery, item); err != nil {
			_ = tx.Rollback()
			return err
		}

		// Update product quantity
		updateQuery := `
			UPDATE products 
			SET quantity = quantity - $1, 
				updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
		`
		if _, err := tx.Exec(updateQuery, item.Quantity, item.ProductId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *DeliveryRepository) DeleteDelivery(deliveryId uuid.UUID) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Get delivery items to restore quantities
	items, err := repo.getDeliveryItems(deliveryId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Restore product quantities
	for _, item := range items {
		updateQuery := `
			UPDATE products 
			SET quantity = quantity + $1, 
				updated_at = CURRENT_TIMESTAMP 
			WHERE id = $2
		`
		if _, err := tx.Exec(updateQuery, item.Quantity, item.ProductId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Delete delivery (cascade will handle items)
	if _, err := tx.Exec(`DELETE FROM deliveries WHERE id = $1`, deliveryId); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *DeliveryRepository) getDeliveryItems(deliveryId uuid.UUID) ([]domain.DeliveryItem, error) {
	items := []domain.DeliveryItem{}
	query := `
		SELECT di.*, p.* 
		FROM delivery_items di
		JOIN products p ON di.product_id = p.id
		WHERE di.delivery_id = $1
	`

	rows, err := repo.db.Queryx(query, deliveryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item domain.DeliveryItem
		var product domain.Product

		// Scan into both structs
		err := rows.Scan(
			&item.Id, &item.DeliveryId, &item.ProductId, &item.Quantity,
			&item.CreatedAt, &item.UpdatedAt,
			&product.Id, &product.Name, &product.Description, &product.SKU,
			&product.Code, &product.Quantity, &product.RestockLevel,
			&product.OptimalLevel, &product.Cost, &product.Price,
			&product.InventoryId, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		item.Product = &product
		items = append(items, item)
	}

	return items, nil
}
