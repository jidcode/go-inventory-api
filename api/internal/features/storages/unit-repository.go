package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ventry/internal/domain"
)

func (repo *StorageRepository) ListUnits(storageId uuid.UUID) ([]domain.StorageUnit, error) {
	units := []domain.StorageUnit{}
	query := `SELECT * FROM storage_units s
			  	WHERE s.storage_id = $1
				GROUP BY s.id
				ORDER BY s.name`

	if err := repo.db.Select(&units, query, storageId); err != nil {
		return nil, err
	}

	// For each unit, fetch its items with product information
	for i := range units {
		items, err := repo.ListUnitItems(units[i].Id)
		if err != nil {
			return nil, err
		}
		units[i].Items = items
	}

	return units, nil
}

func (repo *StorageRepository) GetUnit(unitId uuid.UUID) (*domain.StorageUnit, error) {
	var unit domain.StorageUnit
	query := `SELECT * FROM storage_units WHERE id = $1`

	if err := repo.db.Get(&unit, query, unitId); err != nil {
		return nil, err
	}

	items, err := repo.ListUnitItems(unitId)
	if err != nil {
		return nil, err
	}
	unit.Items = items

	return &unit, nil
}

func (repo *StorageRepository) CreateUnit(unit *domain.StorageUnit) error {
	query := `
		INSERT INTO storage_units (
			id, name, label, storage_id, is_occupied, created_at, updated_at
		) VALUES (
			:id, :name, :label, :storage_id, :is_occupied, :created_at, :updated_at
		)
	`

	_, err := repo.db.NamedExec(query, unit)
	return err
}

func (repo *StorageRepository) EditUnit(unit *domain.StorageUnit) error {
	query := `UPDATE storage_units
				SET name = :name, label = :label, is_occupied = :is_occupied, updated_at = :updated_at
				WHERE id = :id`

	_, err := repo.db.NamedExec(query, unit)
	return err
}

func (repo *StorageRepository) DeleteUnit(unitId uuid.UUID) error {
	_, err := repo.db.Exec(`DELETE FROM storage_units WHERE id = $1`, unitId)
	return err
}

// HELPER FUNCTIONS
func GenerateUnitName(number int) string {
	return fmt.Sprintf("Unit-%05d", number)
}

func (repo *StorageRepository) ListUnitItems(unitId uuid.UUID) ([]domain.UnitItem, error) {
	items := []domain.UnitItem{}

	query := `SELECT * FROM unit_items WHERE storage_unit_id = $1`
	if err := repo.db.Select(&items, query, unitId); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *StorageRepository) AddItemToUnit(unitId, productId uuid.UUID, quantity int) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	var unit domain.StorageUnit
	if err := tx.Get(&unit, `SELECT * FROM storage_units WHERE id = $1`, unitId); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Check if there's already an item in this unit
	var existingItem domain.UnitItem
	err = tx.Get(&existingItem, `SELECT * FROM unit_items WHERE storage_unit_id = $1`, unitId)

	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return err
	}

	if err == sql.ErrNoRows {
		// Create new item
		item := domain.UnitItem{
			Id:            uuid.New(),
			StorageUnitId: unitId,
			ProductId:     productId,
			Quantity:      quantity,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		insertQuery := `
			INSERT INTO unit_items (
				id, storage_unit_id, product_id, quantity, created_at, updated_at
			) VALUES (
				:id, :storage_unit_id, :product_id, :quantity, :created_at, :updated_at
			)
		`

		if _, err := tx.NamedExec(insertQuery, item); err != nil {
			_ = tx.Rollback()
			return err
		}

		// Update unit status
		if _, err := tx.Exec(
			`UPDATE storage_units SET is_occupied = true, updated_at = CURRENT_TIMESTAMP WHERE id = $1`,
			unitId,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
	} else {
		// Update existing item if it's the same product
		if existingItem.ProductId != productId {
			_ = tx.Rollback()
			return errors.New("unit already contains a different product")
		}

		updateQuery := `
			UPDATE unit_items 
			SET quantity = quantity + $1, updated_at = CURRENT_TIMESTAMP
			WHERE storage_unit_id = $2
		`

		if _, err := tx.Exec(updateQuery, quantity, unitId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (repo *StorageRepository) RemoveItemFromUnit(unitId uuid.UUID, quantity int) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
		}
	}()

	// Get current item quantity
	var currentQuantity int
	if err := tx.Get(
		&currentQuantity,
		`SELECT quantity FROM unit_items WHERE storage_unit_id = $1`,
		unitId,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	if currentQuantity < quantity {
		_ = tx.Rollback()
		return errors.New("not enough quantity to remove")
	}

	// Update the quantity
	updateQuery := `
		UPDATE unit_items 
		SET quantity = quantity - $1, updated_at = CURRENT_TIMESTAMP
		WHERE storage_unit_id = $2
	`

	if _, err := tx.Exec(updateQuery, quantity, unitId); err != nil {
		_ = tx.Rollback()
		return err
	}

	// Check if the item quantity is now zero and remove it if necessary
	if currentQuantity == quantity {
		deleteQuery := `DELETE FROM unit_items WHERE storage_unit_id = $1`
		if _, err := tx.Exec(deleteQuery, unitId); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
