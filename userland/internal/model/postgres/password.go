package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/google/uuid"
)

type PasswordStore struct {
	db *sql.DB
}

func NewPasswordStore(db *sql.DB) entity.PasswordStore {
	return &PasswordStore{
		db: db,
	}
}

func (s *PasswordStore) AddNewPassword(ctx context.Context, userId, password string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	query := `INSERT INTO "password" (id, user_id, password) VALUES ($1, $2, $3)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id, userId, password)
	if err != nil {
		return err
	}

	return nil
}

func (s *PasswordStore) GetLastThreePassword(ctx context.Context, userId string) ([]string, error) {
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
