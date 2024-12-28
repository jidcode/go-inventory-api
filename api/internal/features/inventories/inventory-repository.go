package inventories

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain"
	"github.com/ventry/internal/pkg/errors"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(data *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: data}
}

func (repo *InventoryRepository) ListInventories(userId uuid.UUID) ([]domain.Inventory, error) {
	inventories := []domain.Inventory{}
	query := `SELECT * FROM inventories WHERE user_id = $1`

	err := repo.db.Select(&inventories, query, userId)
	if err != nil {
		return nil, errors.DatabaseError(err, "List Inventories")
	}

	return inventories, nil
}

func (repo *InventoryRepository) GetInventory(inventoryId uuid.UUID) (domain.Inventory, error) {
	var inventory domain.Inventory
	query := `SELECT * FROM inventories WHERE id = $1`

	err := repo.db.Get(&inventory, query, inventoryId)
	if err != nil {
		return inventory, errors.DatabaseError(err, "Get Inventory")
	}

	return inventory, nil
}

func (repo *InventoryRepository) CreateInventory(newInventory *domain.Inventory) error {
	query := `INSERT 
				INTO inventories(id, name, description, user_id, created_at, updated_at)
				VALUES(:id, :name, :description, :user_id, :created_at, :updated_at)`

	_, err := repo.db.NamedExec(query, newInventory)
	if err != nil {
		return errors.DatabaseError(err, "Create Inventory")
	}

	return nil
}

func (repo *InventoryRepository) EditInventory(updatedInventory *domain.Inventory) error {
	query := `UPDATE inventories
				SET name = :name, description = :description, user_id = :user_id,
					created_at = :created_at, updated_at = :updated_at
				WHERE id = :id`

	_, err := repo.db.NamedExec(query, updatedInventory)
	if err != nil {
		return errors.DatabaseError(err, "Edit Inventory")
	}

	return nil
}

func (repo *InventoryRepository) DeleteInventory(inventoryId uuid.UUID) error {
	query := `DELETE FROM inventories WHERE id = $1`

	_, err := repo.db.Exec(query, inventoryId)
	if err != nil {
		return errors.DatabaseError(err, "Delete Inventory")
	}

	return nil
}
