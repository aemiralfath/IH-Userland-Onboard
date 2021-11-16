package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

type PlaceStore struct {
	db *sql.DB
}

func NewPlaceStore(db *sql.DB) entity.PlaceStore {
	return &PlaceStore{
		db: db,
	}
}

func (s *PlaceStore) CheckPlaceExist(ctx context.Context, id string) (bool, entity.Place, error) {
	var place entity.Place

	query := `SELECT "id", "name", "description", "current_capacity", "max_capacity", "max_hours" FROM "place" WHERE "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, place, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&place.ID, &place.Name, &place.Description, &place.CurrentCapacity, &place.MaxCapacity, &place.MaxHours)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, place, nil
		} else {
			return false, place, err
		}
	}

	return true, place, nil
}

func (s *PlaceStore) AddNewPlace(ctx context.Context, place entity.Place) (entity.Place, error) {
	query := `INSERT INTO "place" (id, name, description, current_capacity, max_capacity, max_hours) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return place, err
	}

	_, err = stmt.ExecContext(ctx, place.ID, &place.Name, &place.Description, 0, 40, 10)
	if err != nil {
		return place, err
	}

	_, place, err = s.CheckPlaceExist(ctx, place.ID)
	if err != nil {
		return place, err
	}

	return place, nil
}

func (s *PlaceStore) UpdateCurrentCapacity(ctx context.Context, id string, currentcapacity int) (entity.Place, error) {
	var place entity.Place

	query := `UPDATE "place" SET "current_capacity" = $1 WHERE "id" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return place, err
	}

	_, err = stmt.ExecContext(ctx, currentcapacity, id)
	if err != nil {
		return place, err
	}

	_, place, err = s.CheckPlaceExist(ctx, id)
	if err != nil {
		return place, err
	}

	return place, nil
}
