package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (rs *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rs.StatusCode)
	return nil
}

// Success response with (optional) body.
func SuccesRenderer(message string) *Response {
	return &Response{
		StatusCode: 200,
		Message:    message,
	}
}

// HTTP status code if server cannot parse request body or query params. It possibly caused by mismatch data type
func BadRequestErrorRenderer(err error) *Response {
	return &Response{
		Err:        err,
		StatusCode: 400,
		Message:    err.Error(),
	}
}

// HTTP status code for request without authorization or expired token
func UnauthorizedErrorRenderer(err error) *Response {
	return &Response{
		Err:        err,
		StatusCode: 401,
		Message:    err.Error(),
	}
}

// HTTP status code if request sent with authorized token but trying to access to a resource outside its scope.
func ForbiddenErrorRenderer(err error) *Response {
	return &Response{
		Err:        err,
		StatusCode: 403,
		Message:    err.Error(),
	}
}

// Special HTTP status code for input validation error
func UnprocessableErrorRenderer(err error) *Response {
	return &Response{
		Err:        err,
		StatusCode: 403,
		Message:    err.Error(),
	}
}

// Server cannot handle internal error and will not return actual internal error message.
func InternalServerErrorRenderer(err error) *Response {
	return &Response{
		Err:        err,
		StatusCode: 500,
		Message:    err.Error(),
	}
}
