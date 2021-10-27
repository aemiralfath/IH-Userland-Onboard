package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type EventStore struct {
	db *sql.DB
}

func NewEventStore(db *sql.DB) datastore.EventStore {
	return &EventStore{
		db: db,
	}
}

func (s *EventStore) GetEventBySession(ctx context.Context, sessionId string) (*datastore.Event, error) {
	return nil, nil
}

func (s *EventStore) AddNewEvent(ctx context.Context, sessionId string, clientId float64) error {
	query := `INSERT INTO "event" (session_id, client_id) VALUES ($1, $2)`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, sessionId, clientId)
	if err != nil {
		return err
	}

	return nil
}
