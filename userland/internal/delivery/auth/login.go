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

func (d *DeliveryAuth) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := security.GetIPAddress(r)
	clientName := r.Header.Get("X-API-ClientID")
	var body model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Write(w, http.StatusBadRequest, "Invalid body request.", nil, "AUTH-DLV-01")
		return
	}

	if strings.TrimSpace(body.Email) == "" {
		response.Write(w, http.StatusBadRequest, "Email cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.Password) == "" {
		response.Write(w, http.StatusBadRequest, "Password cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	res, err := d.auth.Login(ctx, ip, clientName, body)
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
