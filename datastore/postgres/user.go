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

func (us *UserStore) AddNewUser(ctx context.Context, u *models.User, profile *models.Profile, pass *models.Password) error {

	sql := `INSERT INTO "user" (email, password) VALUES ($1, $2) RETURNING id`
	stmt, err := us.db.Prepare(sql)
	if err != nil {
		return err
	}

	var userId int
	err = stmt.QueryRowContext(ctx, u.Email, u.Password).Scan(&userId)
	if err != nil {
		return err
	}

	_, err = us.db.ExecContext(ctx, `INSERT INTO "profile" (user_id, fullname) VALUES ($1, $2)`, userId, profile.Fullname)
	if err != nil {
		return err
	}

	_, err = us.db.ExecContext(ctx, `INSERT INTO "password" (user_id, password) VALUES ($1, $2)`, userId, pass.Password)
	if err != nil {
		return err
	}
	
	return nil
}
