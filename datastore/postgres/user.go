package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/models"
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

func (us *UserStore) AddNewUser(ctx context.Context, u *models.User) error {
	query := `INSERT INTO "user" (email, password, verified) VALUES ($1, $2, $3)`
	_, err := us.db.ExecContext(ctx, query, u.Email, u.Password, false)
	if err != nil {
		return err
	}
	return nil
}
