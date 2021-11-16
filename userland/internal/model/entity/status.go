package entity

import (
	"context"
	"database/sql"
)

type Report struct {
	Status  Status  `json:"status"`
	Profile Profile `json:"user"`
}

type Status struct {
	ID         string       `json:"id"`
	ProfileId  string       `json:"user_id" validate:"required"`
	PlaceId    string       `json:"place_id" validate:"required"`
	IsGuest    bool         `json:"is_guest"`
	CheckInAt  sql.NullTime `json:"checkin_at"`
	CheckOutAt sql.NullTime `json:"checkout_at"`
	UpdatedAt  sql.NullTime `json:"updated_at"`
}

type StatusStore interface {
	CheckStatusExistById(ctx context.Context, id string) (bool, Status, error)
	CheckStatusExist(ctx context.Context, profileId string, placeId string) (bool, Status, error)
	AddNewStatus(ctx context.Context, profileId, placeId string) (Status, error)
	UpdateCheckOut(ctx context.Context, id string) (Status, error)
	GetStatus(ctx context.Context, id string) (Status, error)
	GetTodayStatus(ctx context.Context, placeId string) ([]Report, error)
}
