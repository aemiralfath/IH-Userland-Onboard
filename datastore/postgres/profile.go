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

func (us *ProfileStore) GetProfile(ctx context.Context, userId float64) (*datastore.Profile, error) {
	sql := `SELECT * FROM "profile" WHERE "user_id" = $1`
	stmt, err := us.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	var prof datastore.Profile
	err = stmt.QueryRowContext(ctx, userId).Scan(&prof.ID, &prof.UserId, &prof.Fullname, &prof.Location, &prof.Bio, &prof.Web, &prof.Picture, &prof.CreatedAt, &prof.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &prof, nil
}
