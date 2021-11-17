package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/entity"
	"github.com/google/uuid"
)

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) entity.SessionStore {
	return &SessionStore{
		db: db,
	}
}

func (s *SessionStore) AddNewSession(ctx context.Context, userId, ip, client string) (string, error) {
	jti, err := uuid.NewRandom()
	if err != nil {
		return jti.String(), err
	}

	query := `INSERT INTO "session" (jti, user_id, ip, client) VALUES ($1, $2, $3, $4)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return jti.String(), err
	}

	_, err = stmt.ExecContext(ctx, jti, userId, ip, client)
	if err != nil {
		return jti.String(), err
	}

	return jti.String(), nil
}

func (s *SessionStore) GetUserSession(ctx context.Context, userId string) ([]entity.Session, error) {
	sessions := []entity.Session{}
	query := `SELECT "jti", "ip", "client", "created_at" FROM "session" WHERE "user_id" = $1`
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
		var session entity.Session
		err := res.Scan(&session.JTI, &session.IP, &session.Client, &session.CreatedAt)
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

func (s *SessionStore) EndSession(ctx context.Context, jti string) error {
	query := `DELETE from "session" WHERE "jti" = $1`
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
