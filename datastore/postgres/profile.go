package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type ProfileStore struct {
	db *sql.DB
}

func NewProfileStore(db *sql.DB) datastore.ProfileStore {
	return &ProfileStore{
		db: db,
	}
}

func (s *ProfileStore) GetProfile(ctx context.Context, userId float64) (*datastore.Profile, error) {
	var prof datastore.Profile
	query := `SELECT * FROM "profile" WHERE "user_id" = $1`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, userId).Scan(
		&prof.ID,
		&prof.UserId,
		&prof.Fullname,
		&prof.Location,
		&prof.Bio,
		&prof.Web,
		&prof.Picture,
		&prof.CreatedAt,
		&prof.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &prof, nil
}

func (s *ProfileStore) AddNewProfile(ctx context.Context, profile *datastore.Profile, userId float64) error {
	query := `INSERT INTO "profile" (user_id, fullname) VALUES ($1, $2)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, profile.Fullname)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProfileStore) UpdateProfile(ctx context.Context, profile *datastore.Profile, userId float64) error {
	query := `UPDATE "profile" SET "fullname" = $1, "location" = $2, "bio" = $3, "web" = $4, "updated_at" = $5 WHERE "user_id" = $6`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, profile.Fullname, profile.Location, profile.Bio, profile.Web, time.Now().Format(time.RFC3339), userId)
	if err != nil {
		return err
	}

	return nil
}

func (s *ProfileStore) UpdatePicture(ctx context.Context, profile *datastore.Profile, userId float64) error {
	query := `UPDATE "profile" SET "picture" = $1 WHERE "user_id" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, profile.Picture, userId)
	if err != nil {
		return err
	}

	return nil
}
