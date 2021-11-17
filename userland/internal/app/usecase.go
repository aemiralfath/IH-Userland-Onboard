package app

import (
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/usecase/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/usecase/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/usecase/session"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/usecase/status"
	"github.com/go-redis/redis/v8"
)

type usecase struct {
	status  *status.UsecaseStatus
	auth    *auth.UsecaseAuth
	session *session.UsecaseSession
	me      *me.UsecaseMe
}

func initUseCase(db *sql.DB, redis *redis.Client) (*usecase, error) {
	r := initRepository(db, redis)
	auth, err := auth.New(r.auth)
	if err != nil {
		return &usecase{}, err
	}

	return &usecase{
		status:  status.New(r.status),
		auth:    auth,
		session: session.New(r.session),
		me:      me.New(r.me),
	}, nil
}
