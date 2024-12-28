package sales

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
)

type SaleRepository struct {
	db *sqlx.DB
}

func NewSaleRepository(db *sqlx.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (repo *SaleRepository) ListSales(inventoryId uuid.UUID) ([]domain.Sale, error) {
	sales := []domain.Sale{}
	query := `
		SELECT * FROM sales s
		WHERE s.inventory_id = $1
		GROUP BY s.id
		ORDER BY s.created_at DESC
	`

	if err := repo.db.Select(&sales, query, inventoryId); err != nil {
		return nil, err
	}

	// For each sale, fetch its items
	for i := range sales {
		items, err := repo.ListSaleItems(sales[i].Id)
		if err != nil {
			return nil, err
		}
		sales[i].Items = items
	}

	return sales, nil
}

func (repo *SaleRepository) GetSale(saleId uuid.UUID) (*domain.Sale, error) {
	var sale domain.Sale
	query := `SELECT * FROM sales WHERE id = $1`

	if err := repo.db.Get(&sale, query, saleId); err != nil {
		return nil, err
	}

	items, err := repo.ListSaleItems(saleId)
	if err != nil {
		return nil, err
	}
	sale.Items = items

	return &sale, nil
}

func (repo *SaleRepository) CreateSale(sale *domain.Sale) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Insert the sale first
	query := `
        INSERT INTO sales (
            id, inventory_id, customer_name, customer_contact, 
            total_amount, created_at, updated_at
        ) VALUES (
            :id, :inventory_id, :customer_name, :customer_contact,
            :total_amount, :created_at, :updated_at
        )
    `

	if _, err := tx.NamedExec(query, sale); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert all sale items
	for i := range sale.Items {
		item := &sale.Items[i]
		item.Id = uuid.New()
		item.SaleId = sale.Id
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
		item.Subtotal = float64(item.Quantity) * item.UnitPrice

		itemQuery := `
            INSERT INTO sale_items (
                id, sale_id, product_id, quantity,
                unit_price, subtotal, created_at, updated_at
            ) VALUES (
                :id, :sale_id, :product_id, :quantity,
                :unit_price, :subtotal, :created_at, :updated_at
            )
        `

		if _, err := tx.NamedExec(itemQuery, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	// Update the total amount based on all items
	updateQuery := `
        UPDATE sales
        SET total_amount = (
            SELECT COALESCE(SUM(subtotal), 0)
            FROM sale_items
            WHERE sale_id = $1
        )
        WHERE id = $1
    `

	if _, err := tx.Exec(updateQuery, sale.Id); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *SaleRepository) EditSale(sale *domain.Sale) error {
	query := `
		UPDATE sales
		SET inventory_id = :inventory_id,
			customer_name = :customer_name,
			customer_contact = :customer_contact,
			total_amount = :total_amount,
			updated_at = :updated_at
		WHERE id = :id
	`

	_, err := repo.db.NamedExec(query, sale)
	return err
}

func (repo *SaleRepository) DeleteSale(saleId uuid.UUID) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.Exec(`DELETE FROM sales WHERE id = $1`, saleId); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// HELPERS
func (repo *SaleRepository) ListSaleItems(saleId uuid.UUID) ([]domain.SaleItem, error) {
	items := []domain.SaleItem{}
	query := `
		SELECT * FROM sale_items
		WHERE sale_id = $1
		ORDER BY created_at
	`

	if err := repo.db.Select(&items, query, saleId); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *SaleRepository) AddItemToSale(item *domain.SaleItem) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Calculate subtotal
	item.Subtotal = float64(item.Quantity) * item.UnitPrice

	// Insert the sale item
	itemQuery := `
		INSERT INTO sale_items (
			id, sale_id, product_id, quantity,
			unit_price, subtotal, created_at, updated_at
		) VALUES (
			:id, :sale_id, :product_id, :quantity,
			:unit_price, :subtotal, :created_at, :updated_at
		)
	`

	if _, err := tx.NamedExec(itemQuery, item); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Update the total amount in the sale
	updateQuery := `
		UPDATE sales
		SET total_amount = (
			SELECT COALESCE(SUM(subtotal), 0)
			FROM sale_items
			WHERE sale_id = $1
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if _, err := tx.Exec(updateQuery, item.SaleId); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (repo *SaleRepository) RemoveItemFromSale(saleId, itemId uuid.UUID) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Delete the sale item
	result, err := tx.Exec(`DELETE FROM sale_items WHERE id = $1 AND sale_id = $2`, itemId, saleId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if rows == 0 {
		_ = tx.Rollback()
		return errors.New("sale item not found")
	}

	// Update the total amount in the sale
	updateQuery := `
		UPDATE sales
		SET total_amount = (
			SELECT COALESCE(SUM(subtotal), 0)
			FROM sale_items
			WHERE sale_id = $1
		),
		updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if _, err := tx.Exec(updateQuery, saleId); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
