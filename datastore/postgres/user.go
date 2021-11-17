package postgres

import (
	"context"
	"database/sql"
	"time"

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

func (s *UserStore) GetUserByEmail(ctx context.Context, email string) (*datastore.User, error) {
	query := `SELECT "id", "email", "password" FROM "user" WHERE "deleted_at" IS NULL AND "verified" = true AND "email" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var user datastore.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserStore) AddNewUser(ctx context.Context, u *datastore.User) (float64, error) {
	var userId float64
	query := `INSERT INTO "user" (email, password) VALUES ($1, $2) RETURNING id`

	stmt, err := us.db.Prepare(query)
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRowContext(ctx, u.Email, u.Password).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *UserStore) ChangePassword(ctx context.Context, u *datastore.User) error {
	query := `UPDATE "user" SET "password" = $1 WHERE "id" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, u.Password, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetEmailByID(ctx context.Context, userId float64) (string, error) {
	query := `SELECT "email" FROM "user" WHERE "deleted_at" IS NULL AND "verified" = true AND "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return "", err
	}

	var email string
	err = stmt.QueryRowContext(ctx, userId).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (s *UserStore) CheckUserEmailExist(ctx context.Context, email string) (*datastore.User, error) {
	query := `SELECT "id", "email" FROM "user" WHERE "deleted_at" IS NULL AND "email" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var user datastore.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) SoftDeleteUser(ctx context.Context, email string) error {
	query := `UPDATE "user" SET "deleted_at" = $1 WHERE "email" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, time.Now().Format(time.RFC3339), email)
	if err != nil {
		return err
	}

	return nil
}
