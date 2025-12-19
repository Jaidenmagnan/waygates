package repositories

import (
	"database/sql"

	"github.com/Jaidenmagnan/waygates/models"
)

type WaygateRepository struct {
	db *sql.DB
}

func NewWaygateRepository(db *sql.DB) *WaygateRepository {
	return &WaygateRepository{db: db}
}

// Create a new waygate in the database.
func (r *WaygateRepository) Create(waygate models.CreateWaygate) (models.Waygate, error) {
	query := "INSERT INTO waygates (name, user_id) VALUES (?, ?)"

	result, err := r.db.Exec(query, waygate.Name, waygate.UserId)
	if err != nil {
		return models.Waygate{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Waygate{}, err
	}

	return models.Waygate{
		ID:     int(id),
		Name:   waygate.Name,
		UserId: waygate.UserId,
	}, nil
}

// Retrieves a waygate by ID.
func (r *WaygateRepository) GetByID(id int) (models.Waygate, error) {
	query := "SELECT id, name, user_id FROM waygates WHERE id = ?"

	row := r.db.QueryRow(query, id)

	var waygate models.Waygate
	if err := row.Scan(&waygate.ID, &waygate.Name, &waygate.UserId); err != nil {
		return models.Waygate{}, err
	}
	return waygate, nil
}

// Deletes a waygate by ID.
func (r *WaygateRepository) DeleteByID(id int) error {
	query := "DELETE FROM waygates WHERE id = ?"

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Update a waygate.
func (r *WaygateRepository) Update(waygate models.Waygate) error {
	query := "UPDATE waygates SET name = ? WHERE id = ?"

	_, err := r.db.Exec(query, waygate.Name, waygate.ID)
	if err != nil {
		return err
	}
	return nil
}

// Get all waygates for a given user.
func (r *WaygateRepository) GetByUserID(userId int) ([]models.Waygate, error) {
	query := "SELECT id, name, user_id FROM waygates WHERE user_id = ?"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waygates []models.Waygate
	for rows.Next() {
		var waygate models.Waygate
		if err := rows.Scan(&waygate.ID, &waygate.Name, &waygate.UserId); err != nil {
			return nil, err
		}
		waygates = append(waygates, waygate)
	}
	return waygates, nil
}
