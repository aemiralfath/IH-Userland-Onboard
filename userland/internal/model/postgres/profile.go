package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/google/uuid"
)

type ProfileStore struct {
	db *sql.DB
}

func NewProfileStore(db *sql.DB) entity.ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

func (s *ProfileStore) CheckNIKExist(ctx context.Context, nik string) (bool, entity.Profile, error) {
	var profile entity.Profile

	query := `SELECT "id", "nik", "name", "status_1", "status_2", "dosage_type", "created_at"::timestamptz(0) FROM "profile" WHERE "nik" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, profile, err
	}

	err = stmt.QueryRowContext(ctx, nik).Scan(&profile.ID, &profile.NIK, &profile.Name, &profile.Status1, &profile.Status2, &profile.DosageType, &profile.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, profile, nil
		} else {
			return false, profile, err
		}
	}

	return true, profile, nil
}

func (s *ProfileStore) GetProfileById(ctx context.Context, id string) (entity.Profile, error) {
	var profile entity.Profile

	query := `SELECT "id", "nik", "name", "status_1", "status_2", "dosage_type", "picture", "created_at"::timestamptz(0), "updated_at"::timestamptz(0)  FROM "profile" WHERE "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return profile, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&profile.ID, &profile.NIK, &profile.Name, &profile.Status1, &profile.Status2, &profile.DosageType, &profile.Picture, &profile.CreatedAt, &profile.UpdatedAt)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s *ProfileStore) AddNewProfile(ctx context.Context, profile entity.Profile) (entity.Profile, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return profile, err
	}

	query := `INSERT INTO "profile" (id, nik, name, status_1, status_2, dosage_type) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return profile, err
	}

	_, err = stmt.ExecContext(ctx, id, &profile.NIK, &profile.Name, &profile.Status1, &profile.Status2, &profile.DosageType)
	if err != nil {
		return profile, err
	}

	_, profile, err = s.CheckNIKExist(ctx, profile.NIK)
	if err != nil {
		return profile, err
	}

	return profile, nil
}

func (s *ProfileStore) UpdateProfile(ctx context.Context, id, name, dosageType string) error {
	query := `UPDATE "profile" SET "name" = $1, "dosage_type" = $2, "updated_at" = $3 WHERE "id" = $4`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, name, dosageType, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return err
	}

	return nil
}
