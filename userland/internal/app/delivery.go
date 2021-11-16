package app

import (
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/delivery/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/delivery/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/delivery/session"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/delivery/status"
	"github.com/go-redis/redis/v8"
)

type delivery struct {
	status  *status.DeliveryStatus
	auth    *auth.DeliveryAuth
	session *session.DeliverySession
	me      *me.DeliveryMe
}

func initDelivery(db *sql.DB, redis *redis.Client) delivery {
	u := initUseCase(db, redis)
	return delivery{
		status:  status.NewStatus(u.status),
		auth:    auth.NewAuth(u.auth),
		session: session.NewSession(u.session),
		me:      me.NewMe(u.me),
	}
}
