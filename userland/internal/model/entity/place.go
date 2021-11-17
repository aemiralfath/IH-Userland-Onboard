package entity

import "context"

type Place struct {
	ID              string `json:"id"`
	Name            string `json:"name" validate:"required"`
	Description     string `json:"description" validate:"required"`
	CurrentCapacity int    `json:"current_capacity"`
	MaxCapacity     int    `json:"max_capacity"`
	MaxHours        int    `json:"max_hours"`
}

type PlaceStore interface {
	CheckPlaceExist(ctx context.Context, id string) (bool, Place, error)
	AddNewPlace(ctx context.Context, place Place) (Place, error)
	UpdateCurrentCapacity(ctx context.Context, id string, currentcapacity int) (Place, error)
}
