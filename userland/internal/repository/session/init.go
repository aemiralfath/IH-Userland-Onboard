package session

import (
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/postgres"
)

type Repository struct {
	UserStore     entity.UserStore
	PasswordStore entity.PasswordStore
	SessionStore  entity.SessionStore
	ProfileStore  entity.ProfileStore
	StatusStore   entity.StatusStore
	PlaceStore    entity.PlaceStore
}

func New(db *sql.DB) *Repository {
	return &Repository{
		UserStore:     postgres.NewUserStore(db),
		PasswordStore: postgres.NewPasswordStore(db),
		SessionStore:  postgres.NewSessionStore(db),
		ProfileStore:  postgres.NewProfileStore(db),
		StatusStore:   postgres.NewStatusStore(db),
		PlaceStore:    postgres.NewPlaceStore(db),
	}
}
