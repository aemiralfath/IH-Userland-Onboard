package me

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
)

func (d *DeliveryMe) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body model.UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Write(w, http.StatusBadRequest, "Invalid body request.", nil, "AUTH-DLV-01")
		return
	}

	if strings.TrimSpace(body.Fullname) == "" {
		response.Write(w, http.StatusBadRequest, "Name cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.DosageType) == "" {
		response.Write(w, http.StatusBadRequest, "Dosage Type cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	_, claims, err := jwt.New().FromContext(ctx)
	if err != nil {
		e, ok := err.(*myerror.Error)
		if !ok {
			response.Write(w, http.StatusInternalServerError, "Our server encounter a problem.", nil, "BAD-ERROR")
			return
		}
		response.Write(w, http.StatusBadRequest, e.Error(), nil, e.ErrorCode)
		return
	}

	userId := claims["userId"]
	res, err := d.me.UpdateProfile(ctx, userId.(string), body)
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
