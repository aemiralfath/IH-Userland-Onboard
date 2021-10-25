package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
)

type resetPasswordRequest struct {
	Token           string `json:"token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

func ResetPassword(userStore datastore.UserStore, passwordStore datastore.PasswordStore, otp datastore.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &resetPasswordRequest{}

		if err := render.Bind(r, req); err != nil {
			render.Render(w, r, helper.BadRequestErrorRenderer(err))
			return
		}

		hashPassword, err := helper.HashPassword(req.Password)
		if err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		req.Password = string(hashPassword)
		if err := userStore.ChangePassword(ctx, parseResetRequestUser(req)); err != nil {
			render.Render(w, r, helper.InternalServerErrorRenderer(err))
			return
		}

		success := struct {
			Success bool `json:"success"`
		}{Success: true}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(success)
	}
}

func parseResetRequestUser(p *resetPasswordRequest) *datastore.User {
	return &datastore.User{
		Password: p.Password,
	}
}

func (request *resetPasswordRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(request.Token) == "" {
		return fmt.Errorf("required token")
	}

	if request.Password != request.PasswordConfirm {
		return fmt.Errorf("password and confirm password must same!")
	}

	passLength, number, upper := helper.VerifyPassword(request.Password)
	if !passLength || !number || !upper {
		return fmt.Errorf("password must have lowercase, uppercase, number, and minimum 8 chars!")
	}

	return nil
}

func (*resetPasswordRequest) Render(w http.ResponseWriter, r *http.Request) {}
