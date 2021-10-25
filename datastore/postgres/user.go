package postgres

import (
	"context"
	"database/sql"
	"fmt"

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

func (us *UserStore) GetUser(ctx context.Context, u *datastore.User) (*datastore.User, error) {
	sql := `SELECT "id", "email", "password" FROM "user" WHERE "deleted_at" IS NULL AND "verified" = true AND "email" = $1`
	stmt, err := us.db.Prepare(sql)
	if err != nil {
		return nil, err
	}

	var user datastore.User
	err = stmt.QueryRowContext(ctx, u.Email).Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, fmt.Errorf("User not found")
		} else {
			return nil, err
		}
	} else {
		return &user, nil
	}
}

func (us *UserStore) AddNewUser(ctx context.Context, u *datastore.User) (int, error) {
	var userId int
	sql := `INSERT INTO "user" (email, password) VALUES ($1, $2) RETURNING id`

	stmt, err := us.db.Prepare(sql)
	if err != nil {
		return 0, err
	}

	err = stmt.QueryRowContext(ctx, u.Email, u.Password).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (us *UserStore) ChangePassword(ctx context.Context, u *datastore.User) error {
	return nil
}
