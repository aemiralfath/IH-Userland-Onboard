package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) datastore.SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (s *SessionStore) GetUserSession(ctx context.Context, userId float64) (*datastore.Session, error) {
	return nil, nil
}

func (s *SessionStore) AddNewSession(ctx context.Context, session *datastore.Session, clientId float64) error {
	query := `INSERT INTO "session" (jti, user_id, client_id, is_current) VALUES ($1, $2, $3, $4)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, session.JTI, session.UserId, clientId, session.IsCurrent)
	if err != nil {
		return err
	}

	return nil
}
