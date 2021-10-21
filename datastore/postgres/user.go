package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) datastore.UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) GetUser(ctx context.Context) error {
	_, _ = us.db.QueryContext(ctx, "")
	return nil
}
