package repositories

import (
	"database/sql"
	"sort"

	"github.com/Jaidenmagnan/waygates/models"
)

type WaygateLinkRepository struct {
	db *sql.DB
}

func NewWaygateLinkRepository(db *sql.DB) *WaygateLinkRepository {
	return &WaygateLinkRepository{db: db}
}

// Create a new waygate link in the database.
func (r *WaygateLinkRepository) Create(waygateLink models.CreateWaygateLink) (models.WaygateLink, error) {
	query := "INSERT INTO waygate_links (name, link, waygate_id) VALUES (?, ?, ?)"

	result, err := r.db.Exec(query, waygateLink.Name, waygateLink.Link, waygateLink.WaygateId)
	if err != nil {
		return models.WaygateLink{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.WaygateLink{}, err
	}

	return models.WaygateLink{
		ID:   int(id),
		Name: waygateLink.Name,
		Link: waygateLink.Link,
	}, nil

}

func (r *WaygateLinkRepository) Delete(id int) error {
	query := "DELETE FROM waygate_links WHERE id = ?"

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Get all links for a waygate ID.
func (r *WaygateLinkRepository) GetByWaygateID(waygateID int) ([]models.WaygateLink, error) {
	query := "SELECT id, name, link FROM waygate_links WHERE waygate_id = ?"

	rows, err := r.db.Query(query, waygateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var waygates []models.WaygateLink
	for rows.Next() {
		var waygateLink models.WaygateLink
		if err := rows.Scan(&waygateLink.ID, &waygateLink.Name, &waygateLink.Link); err != nil {
			return nil, err
		}
		waygates = append(waygates, waygateLink)

		sort.Slice(waygates, func(i, j int) bool {
			return waygates[i].Name <= waygates[j].Name
		})
	}
	return waygates, nil
}

// Update a waygate link.
func (r *WaygateLinkRepository) Update(waygateLink models.WaygateLink) error {
	query := "UPDATE waygate_links SET name, link = (?, ?) WHERE id = ?"

	_, err := r.db.Exec(query, waygateLink.Name, waygateLink.Link, waygateLink.ID)
	if err != nil {
		return err
	}
	return nil
}
