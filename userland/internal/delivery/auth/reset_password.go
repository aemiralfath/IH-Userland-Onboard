package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/security"
)

func (d *DeliveryAuth) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body model.ResetPasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Write(w, http.StatusBadRequest, "Invalid body request.", nil, "AUTH-DLV-01")
		return
	}

	if strings.TrimSpace(body.Token) == "" {
		response.Write(w, http.StatusBadRequest, "Token cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if body.Password != body.PasswordConfirm {
		response.Write(w, http.StatusBadRequest, "Password and confirm password must same!.", nil, "AUTH-DLV-02")
		return
	}

	passLength, number, upper := security.VerifyPassword(body.Password)
	if !passLength || !number || !upper {
		response.Write(w, http.StatusBadRequest, "Password must have lowercase, uppercase, number, and minimum 8 chars!", nil, "AUTH-DLV-02")
		return
	}

	res, err := d.auth.ResetPassword(ctx, body)
	if err != nil {
		e, ok := err.(*myerror.Error)
		if !ok {
			response.Write(w, http.StatusInternalServerError, "Our server encounter a problem.", nil, "BAD-ERROR")
			return
		}
		response.Write(w, http.StatusBadRequest, e.Error(), nil, e.ErrorCode)
		return
	}

	response.Write(w, http.StatusOK, "success", res, "")
}
