package app

import (
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/repository/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/repository/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/repository/session"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/repository/status"
	"github.com/go-redis/redis/v8"
)

type repository struct {
	status  *status.Repository
	auth    *auth.Repository
	session *session.Repository
	me      *me.Repository
}

func initRepository(db *sql.DB, redis *redis.Client) *repository {
	return &repository{
		status:  status.New(db),
		auth:    auth.New(db, redis),
		session: session.New(db),
		me:      me.New(db, redis),
	}
}
