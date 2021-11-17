package myerror

import "errors"

type Error struct {
	ErrorCode string
	Err       error
}

func New(message string, errCode string) error {
	return &Error{
		ErrorCode: errCode,
		Err:       errors.New(message),
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
