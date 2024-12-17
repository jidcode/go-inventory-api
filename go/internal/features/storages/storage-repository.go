package storages

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type StorageRepository struct {
	db *sqlx.DB
}

func NewStorageRepository(data *sqlx.DB) *StorageRepository {
	return &StorageRepository{db: data}
}

func (repo *StorageRepository) ListStorages(inventoryID uuid.UUID) ([]models.Storage, error) {
	storages := []models.Storage{}
	query := `SELECT * FROM storages WHERE inventory_id = $1`
	err := repo.db.Select(&storages, query, inventoryID)
	return storages, err
}

func (repo *StorageRepository) GetStorage(id uuid.UUID) (models.Storage, error) {
	var storage models.Storage

	query := `SELECT * FROM storages WHERE id = $1`
	err := repo.db.Get(&storage, query, id)
	if err != nil {
		return storage, err
	}

	unitsQuery := `SELECT * FROM units WHERE storage_id = $1`
	err = repo.db.Select(&storage.Units, unitsQuery, id)
	if err != nil {
		return storage, err
	}

	return storage, nil
}

func (repo *StorageRepository) CreateStorage(storage *models.Storage) error {
	storage.ID = uuid.New()
	storage.CreatedAt = time.Now()
	storage.UpdatedAt = time.Now()

	query := `INSERT INTO storages(
		id, name, location, capacity, occupied_space, inventory_id, created_at, updated_at
	) VALUES (
		:id, :name, :location, :capacity, :occupied_space, :inventory_id, :created_at, :updated_at
	)`

	_, err := repo.db.NamedExec(query, storage)
	return err
}

func (repo *StorageRepository) UpdateStorage(storage *models.Storage) error {
	storage.UpdatedAt = time.Now()
	query := `UPDATE storages
		SET name = :name, 
			location = :location, 
			capacity = :capacity, 
			occupied_space = :occupied_space, 
			inventory_id = :inventory_id, 
			updated_at = :updated_at
		WHERE id = :id`

	_, err := repo.db.NamedExec(query, storage)
	return err
}

func (repo *StorageRepository) DeleteStorage(id uuid.UUID) error {
	query := `DELETE FROM storages WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
