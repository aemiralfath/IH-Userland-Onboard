package session

import (
	"context"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
)

func (r *Repository) ListSession(ctx context.Context, userId string) ([]entity.Session, error) {
	sessions, err := r.SessionStore.GetUserSession(ctx, userId)
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}
