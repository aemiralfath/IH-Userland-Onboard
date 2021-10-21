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

func (us *ProfileStore) GetProfile(ctx context.Context) error {
	_, _ = us.db.QueryContext(ctx, "")
	return nil
}
