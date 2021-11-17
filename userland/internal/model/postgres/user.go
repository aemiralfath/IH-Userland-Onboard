package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/google/uuid"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) entity.UserStore {
	return &UserStore{
		db: db,
	}
}

func (s *UserStore) CheckEmailExist(ctx context.Context, email string) (bool, entity.User, error) {
	var user entity.User
	query := `SELECT "id", "email", "password" FROM "user" WHERE "deleted_at" IS NULL AND "email" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, user, err
	}

	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, user, nil
		} else {
			return false, user, err
		}
	}

	return true, user, nil
}

func (s *UserStore) GetUserById(ctx context.Context, id string) (entity.User, error) {
	var user entity.User
	query := `SELECT "id", "profile_id", "email", "password" FROM "user" WHERE "deleted_at" IS NULL AND "id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return user, err
	}

	err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.ProfileId, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (us *UserStore) AddNewUser(ctx context.Context, profileId, email, password string) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return id.String(), err
	}

	query := `INSERT INTO "user" (id, profile_id, email, password, verified) VALUES ($1, $2, $3, $4, $5)`

	stmt, err := us.db.Prepare(query)
	if err != nil {
		return id.String(), err
	}

	_, err = stmt.ExecContext(ctx, id, profileId, email, password, false)
	if err != nil {
		return id.String(), err
	}

	return id.String(), nil
}

func (s *UserStore) ChangePassword(ctx context.Context, id, password string) error {
	query := `UPDATE "user" SET "password" = $1 WHERE "id" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, password, id)
	if err != nil {
		return err
	}

	return nil
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
