package models

type Profile struct {
	ID        int64
	UserId    int64
	Fullname  string
	Location  string
	Bio       string
	Web       string
	Picture   string
	UpdatedAt string
}