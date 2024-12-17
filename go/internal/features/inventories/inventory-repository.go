package inventories

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type InventoryRepository struct {
	db *sqlx.DB
}

func NewInventoryRepository(data *sqlx.DB) *InventoryRepository {
	return &InventoryRepository{db: data}
}

// List all inventories for a specific user
func (repo *InventoryRepository) ListInventories(userID uuid.UUID) ([]models.Inventory, error) {
	inventories := []models.Inventory{}
	query := `SELECT * FROM inventories WHERE user_id = $1`

	err := repo.db.Select(&inventories, query, userID)
	return inventories, err
}

// Get a single inventory by ID
func (repo *InventoryRepository) GetInventory(id uuid.UUID) (models.Inventory, error) {
	var inventory models.Inventory
	query := `SELECT * FROM inventories WHERE id = $1`
	err := repo.db.Get(&inventory, query, id)
	return inventory, err
}

// Add a new inventory
func (repo *InventoryRepository) CreateInventory(inventory *models.Inventory) error {
	inventory.ID = uuid.New()
	inventory.CreatedAt = time.Now()
	inventory.UpdatedAt = time.Now()

	query := `INSERT 
				INTO inventories(id, name, description, user_id, created_at, updated_at)
				VALUES(:id, :name, :description, :user_id, :created_at, :updated_at)`

	_, err := repo.db.NamedExec(query, inventory)
	return err
}

// Update a specific inventory
func (repo *InventoryRepository) UpdateInventory(inventory *models.Inventory) error {
	inventory.UpdatedAt = time.Now()
	query := `UPDATE inventories
				SET name = :name, description = :description, user_id = :user_id,
					created_at = :created_at, updated_at = :updated_at
				WHERE id = :id`

	_, err := repo.db.NamedExec(query, inventory)
	return err
}

// Delete an inventory
func (repo *InventoryRepository) DeleteInventory(id uuid.UUID) error {
	query := `DELETE FROM inventories WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
