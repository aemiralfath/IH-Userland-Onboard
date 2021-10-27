package postgres

import (
	"context"
	"database/sql"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

type ClientStore struct {
	db *sql.DB
}

func NewClientStore(db *sql.DB) datastore.ClientStore {
	return &ClientStore{
		db: db,
	}
}

func (s *ClientStore) GetClientByName(ctx context.Context, name string) (*datastore.Client, error) {
	query := `SELECT "id", "name" FROM "client" WHERE "name" = $1`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var client datastore.Client
	err = stmt.QueryRowContext(ctx, name).Scan(&client.ID, &client.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			newClient, err := s.AddNewClient(ctx, name)
			if err != nil {
				return nil, err
			}
			return newClient, nil
		} else {
			return nil, err
		}
	} else {
		return &client, nil
	}
}

func (s *ClientStore) AddNewClient(ctx context.Context, name string) (*datastore.Client, error) {
	var clientId float64
	query := `INSERT INTO "client" (name) VALUES ($1) RETURNING id`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	err = stmt.QueryRowContext(ctx, name).Scan(&clientId)
	if err != nil {
		return nil, err
	}

	return &datastore.Client{ID: clientId, Name: name}, nil
}
