package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type PasswordStore struct {
	db *sql.DB
}

func NewPasswordStore(db *sql.DB) datastore.PasswordStore {
	return &PasswordStore{
		db: db,
	}
}

func (us *PasswordStore) GetPassword(ctx context.Context) error {
	_, _ = us.db.QueryContext(ctx, "")
	return nil
}
