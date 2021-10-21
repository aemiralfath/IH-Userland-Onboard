package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/models"
	"github.com/go-chi/render"
)

type registerRequest struct {
	Fullname        string `json:"fullname"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func Register(userStore datastore.UserStore, profileStore datastore.ProfileStore, passwordStore datastore.PasswordStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &registerRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, handler.BadRequestErrorRenderer(err))
			return
		}

		if err := userStore.AddNewUser(ctx, parseHandler(req)); err != nil {
			render.Render(w, r, handler.InternalServerErrorRenderer(err))
			return
		}

		if err := render.Render(w, r, handler.SuccesRenderer("Success")); err != nil {
			render.Render(w, r, handler.InternalServerErrorRenderer(err))
			return
		}
	}
}

func parseHandler(u *registerRequest) *models.User {
	return &models.User{
		Email: u.Email,
		Password: u.Password,
	}
}

func (register *registerRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(register.Fullname) == "" {
		return fmt.Errorf("required fullname")
	}
	return nil
}

func (*registerRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
