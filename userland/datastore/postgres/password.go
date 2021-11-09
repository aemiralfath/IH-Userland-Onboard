package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
)

type PasswordStore struct {
	db *sql.DB
}

func NewPasswordStore(db *sql.DB) datastore.PasswordStore {
	return &PasswordStore{
		db: db,
	}
}

func (s *PasswordStore) GetPassword(ctx context.Context) error {
	_, _ = s.db.QueryContext(ctx, "")
	return nil
}

func (s *PasswordStore) AddNewPassword(ctx context.Context, password *datastore.Password, userId float64) error {
	query := `INSERT INTO "password" (user_id, password) VALUES ($1, $2)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, password.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *PasswordStore) GetLastThreePassword(ctx context.Context, userId float64) ([]string, error) {
	passwords := []string{}
	query := `SELECT "password" FROM "password" WHERE "user_id" = $1 ORDER BY "created_at" DESC LIMIT 3`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return passwords, err
	}

	res, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return passwords, err
	}
	defer res.Close()

	for res.Next() {
		var password string
		err := res.Scan(&password)
		if err != nil {
			return passwords, err
		}

		passwords = append(passwords, password)
	}

	if err = res.Err(); err != nil {
		return passwords, err
	}

	return passwords, nil
}
