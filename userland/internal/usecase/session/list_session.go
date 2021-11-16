package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/rs/zerolog/log"
)

func (u *UsecaseSession) ListSession(ctx context.Context, userId string) (model.ListSessionResponse, error) {
	var result model.ListSessionResponse

	sessions, err := u.session.ListSession(ctx, userId)
	if err != nil {
		log.Error().Err(err).Stack().Msg("Error get sessions")
		return result, myerror.New(err.Error(), "STATUS-USC-01")
	}

	for _, e := range sessions {
		var session model.SessionResponse

		session.JTI = e.JTI
		session.Client = e.Client
		session.IP = e.IP
		session.CreatedAt = e.CreatedAt

		result.Sessions = append(result.Sessions, session)
	}

	return result, nil
}
