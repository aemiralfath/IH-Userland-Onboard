package entity

import (
	"context"
)

type Profile struct {
	ID         string `json:"id"`
	NIK        string `json:"nik" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Status1    bool   `json:"status_1"`
	Status2    bool   `json:"status_2"`
	DosageType string `json:"dosage_type"`
	Picture    string `json:"picture"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ProfileStore interface {
	CheckNIKExist(ctx context.Context, nik string) (bool, Profile, error)
	AddNewProfile(ctx context.Context, profile Profile) (Profile, error)
	GetProfileById(ctx context.Context, id string) (Profile, error)
	UpdateProfile(ctx context.Context, id, name, dosageType string) error
}
