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

func (s *PasswordStore) AddNewPassword(ctx context.Context, password *datastore.Password, userId int) error {
	sql := `INSERT INTO "password" (user_id, password) VALUES ($1, $2)`
	stmt, err := s.db.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, password.Password)
	if err != nil {
		return err
	}

	return nil
}
