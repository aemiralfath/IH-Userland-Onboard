package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/google/uuid"
)

type StatusStore struct {
	db *sql.DB
}

func NewStatusStore(db *sql.DB) entity.StatusStore {
	return &StatusStore{
		db: db,
	}
}

func (s *StatusStore) CheckStatusExist(ctx context.Context, profileId string, placeId string) (bool, entity.Status, error) {
	var status entity.Status

	query := `SELECT "id", "profile_id", "place_id", "checkin_at"::timestamptz(0), "checkout_at"::timestamptz(0), "updated_at"::timestamptz(0) FROM "status" WHERE "profile_id" = $1 AND "place_id" = $2 AND DATE("checkin_at") <= CURRENT_DATE AND "checkout_at" IS NULL`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, status, err
	}

	err = stmt.QueryRowContext(ctx, profileId, placeId).Scan(&status.ID, &status.ProfileId, &status.PlaceId, &status.CheckInAt, &status.CheckOutAt, &status.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, status, nil
		} else {
			return false, status, err
		}
	}

	return true, status, nil
}

func (s *StatusStore) CheckStatusExistById(ctx context.Context, id string) (bool, entity.Status, error) {
	var status entity.Status

	query := `SELECT "id", "profile_id", "place_id", "checkin_at"::timestamptz(0), "checkout_at"::timestamptz(0), "updated_at"::timestamptz(0) FROM "status" WHERE "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, status, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&status.ID, &status.ProfileId, &status.PlaceId, &status.CheckInAt, &status.CheckOutAt, &status.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, status, nil
		} else {
			return false, status, err
		}
	}

	return true, status, nil
}

func (s *StatusStore) AddNewStatus(ctx context.Context, profileId, placeId string) (entity.Status, error) {
	var status entity.Status

	id, err := uuid.NewRandom()
	if err != nil {
		return status, err
	}

	query := `INSERT INTO "status" (id, profile_id, place_id) VALUES ($1, $2, $3)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return status, err
	}

	_, err = stmt.ExecContext(ctx, id, profileId, placeId)
	if err != nil {
		return status, err
	}

	_, status, err = s.CheckStatusExist(ctx, profileId, placeId)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *StatusStore) UpdateCheckOut(ctx context.Context, id string) (entity.Status, error) {
	var status entity.Status

	query := `UPDATE "status" SET "checkout_at" = $1, "updated_at" = $2 WHERE "id" = $3`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return status, err
	}

	_, err = stmt.ExecContext(ctx, time.Now(), time.Now(), id)
	if err != nil {
		return status, err
	}

	status, err = s.GetStatus(ctx, id)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *StatusStore) GetStatus(ctx context.Context, id string) (entity.Status, error) {
	var status entity.Status

	query := `SELECT "id", "profile_id", "place_id", "checkin_at"::timestamptz(0), "checkout_at"::timestamptz(0), "updated_at"::timestamptz(0) FROM "status" WHERE "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return status, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&status.ID, &status.ProfileId, &status.PlaceId, &status.CheckInAt, &status.CheckOutAt, &status.UpdatedAt)
	if err != nil {
		return status, err
	}

	return status, nil
}

func (s *StatusStore) GetTodayStatus(ctx context.Context, placeId string) ([]entity.Report, error) {
	var reports []entity.Report

	query := `SELECT 
				"status"."id", 
				"status"."profile_id", 
				"status"."place_id", 
				"status"."checkin_at"::timestamptz(0), 
				"status"."checkout_at"::timestamptz(0),
				"status"."updated_at"::timestamptz(0),
				"profile"."id",
				"profile"."nik",
				"profile"."name",
				"profile"."status_1",
				"profile"."status_2",
				"profile"."dosage_type",
				"profile"."created_at"::timestamptz(0)
			FROM "status"
			INNER JOIN "profile" ON "profile"."id" = "status"."profile_id"
			WHERE "status"."place_id" = $1 and DATE("checkin_at") = CURRENT_DATE
			ORDER BY "status"."updated_at" DESC`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return reports, err
	}

	res, err := stmt.QueryContext(ctx, placeId)
	if err != nil {
		return reports, err
	}
	defer res.Close()

	for res.Next() {
		var status entity.Status
		var profile entity.Profile

		err := res.Scan(
			&status.ID, &status.ProfileId, &status.PlaceId, &status.CheckInAt, &status.CheckOutAt, &status.UpdatedAt,
			&profile.ID, &profile.NIK, &profile.Name, &profile.Status1, &profile.Status2, &profile.DosageType, &profile.CreatedAt)

		if err != nil {
			return reports, err
		}

		reports = append(reports, entity.Report{Status: status, Profile: profile})
	}

	if err = res.Err(); err != nil {
		return reports, err
	}

	return reports, nil
}
