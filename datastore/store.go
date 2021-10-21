package datastore

import (
	"context"
)

type UserStore interface {
	GetUser(ctx context.Context) error
}

type AuthStore interface {
	GetAuth(ctx context.Context) error
}