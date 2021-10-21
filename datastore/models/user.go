package models

type User struct {
	ID        int64
	Email     string
	Password  string
	Verified  bool
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}