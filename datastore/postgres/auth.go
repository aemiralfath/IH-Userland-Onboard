package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) datastore.AuthStore {
	return &AuthStore{
		db: db,
	}
}

func (us *AuthStore) GetAuth(ctx context.Context) error {
	_, _ = us.db.QueryContext(ctx, "")
	return nil
}
