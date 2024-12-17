package storages

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ventry/internal/domain/models"
)

type UnitsRepository struct {
	db *sqlx.DB
}

func NewUnitsRepository(data *sqlx.DB) *UnitsRepository {
	return &UnitsRepository{db: data}
}

func (repo *UnitsRepository) ListUnits(storageID uuid.UUID) ([]models.Unit, error) {
	units := []models.Unit{}
	query := `SELECT * FROM units WHERE storage_id = $1`
	err := repo.db.Select(&units, query, storageID)
	return units, err
}

func (repo *UnitsRepository) GetUnit(id uuid.UUID) (models.Unit, error) {
	var unit models.Unit
	query := `SELECT * FROM units WHERE id = $1`
	err := repo.db.Get(&unit, query, id)
	return unit, err
}

func (repo *UnitsRepository) CreateUnit(unit *models.Unit) error {
	unit.ID = uuid.New()
	unit.CreatedAt = time.Now()
	unit.UpdatedAt = time.Now()

	query := `INSERT INTO 
			  units (id, label, storage_id, is_occupied, capacity, created_at, updated_at)
			  VALUES (:id, :label, :storage_id, :is_occupied, :capacity, :created_at, :updated_at)`
	_, err := repo.db.NamedExec(query, unit)
	return err
}

func (repo *UnitsRepository) UpdateUnit(unit *models.Unit) error {
	unit.UpdatedAt = time.Now()
	query := `UPDATE units
              SET label = :label, 
                  storage_id = :storage_id, 
                  is_occupied = :is_occupied, 
                  capacity = :capacity, 
                  updated_at = :updated_at
              WHERE id = :id`
	_, err := repo.db.NamedExec(query, unit)
	return err
}

func (repo *UnitsRepository) DeleteUnit(id uuid.UUID) error {
	query := `DELETE FROM units WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
