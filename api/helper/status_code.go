package helper

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Err        error  `json:"-"`
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

type SuccessResponse struct {
	Success    bool `json:"success"`
	StatusCode int  `json:"-"`
}

func CustomRender(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	content, _ := json.Marshal(data)
	_, _ = w.Write(content)
}

func (rs *ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rs.StatusCode)
	return nil
}

func (rs *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, rs.StatusCode)
	return nil
}

// Success response with (optional) body.
func SuccesRenderer() *SuccessResponse {
	return &SuccessResponse{
		StatusCode: 200,
		Success:    true,
	}
}

// HTTP status code if server cannot parse request body or query params. It possibly caused by mismatch data type
func BadRequestErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 400,
		Message:    err.Error(),
	}
}

// HTTP status code for request without authorization or expired token
func UnauthorizedErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 401,
		Message:    err.Error(),
	}
}

// HTTP status code if request sent with authorized token but trying to access to a resource outside its scope.
func ForbiddenErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 403,
		Message:    err.Error(),
	}
}

// Special HTTP status code for input validation error
func UnprocessableErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 403,
		Message:    err.Error(),
	}
}

// Server cannot handle internal error and will not return actual internal error message.
func InternalServerErrorRenderer(err error) *ErrorResponse {
	return &ErrorResponse{
		Err:        err,
		StatusCode: 500,
		Message:    "Internal Server Error",
	}
}
