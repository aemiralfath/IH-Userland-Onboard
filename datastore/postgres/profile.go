package postgres

import (
	"context"
	"database/sql"

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
	sql := `SELECT * FROM "profile" WHERE "user_id" = $1`

	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, userId).Scan(&prof.ID, &prof.UserId, &prof.Fullname, &prof.Location, &prof.Bio, &prof.Web, &prof.Picture, &prof.CreatedAt, &prof.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}

func (s *ProfileStore) AddNewProfile(ctx context.Context, profile *datastore.Profile, userId int) error {
	sql := `INSERT INTO "profile" (user_id, fullname) VALUES ($1, $2)`
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, profile.Fullname)
	if err != nil {
		return err
	}

	return nil
}
