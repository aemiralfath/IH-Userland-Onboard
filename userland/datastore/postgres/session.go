package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) datastore.SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (s *SessionStore) GetUserSession(ctx context.Context, userId float64) ([]datastore.Session, error) {
	sessions := []datastore.Session{}
	query := `SELECT "client_id", "is_current", "ip", "created_at", "updated_at" FROM "session" WHERE "user_id" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return sessions, err
	}

	res, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return sessions, err
	}
	defer res.Close()

	for res.Next() {
		var session datastore.Session
		err := res.Scan(&session.ClientId, &session.IsCurrent, &session.IP, &session.CreatedAt, &session.UpdatedAt)
		if err != nil {
			return sessions, err
		}

		err = s.db.QueryRowContext(ctx, `SELECT "id", "name" FROM "client" WHERE "id" = $1`, session.ClientId).Scan(&session.Client.ID, &session.Client.Name)
		if err != nil {
			return sessions, err
		}

		sessions = append(sessions, session)
	}

	if err = res.Err(); err != nil {
		return sessions, err
	}

	return sessions, nil
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

func (s *SessionStore) EndSession(ctx context.Context, jti string) error {
	query := `UPDATE "session" SET "is_current" = $1 WHERE "jti" = $2`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, false, jti)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionStore) DeleteOtherSession(ctx context.Context, jti string) error {
	query := `DELETE from "session" WHERE "jti" != $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, jti)
	if err != nil {
		return err
	}

	return nil
}
