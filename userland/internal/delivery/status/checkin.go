package status

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
)

func (d *DeliveryStatus) CheckIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var body model.CheckInRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.Write(w, http.StatusBadRequest, "Invalid body request.", nil, "AUTH-DLV-01")
		return
	}

	if strings.TrimSpace(body.Profile.NIK) == "" {
		response.Write(w, http.StatusBadRequest, "Profile NIK cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.Profile.Name) == "" {
		response.Write(w, http.StatusBadRequest, "Profile Name cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if !body.Profile.Status1 || !body.Profile.Status2 {
		response.Write(w, http.StatusBadRequest, "Must be vaccinated 2 doses.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.Place.ID) == "" {
		response.Write(w, http.StatusBadRequest, "Place ID cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.Place.Name) == "" {
		response.Write(w, http.StatusBadRequest, "Place Name cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	if strings.TrimSpace(body.Place.Description) == "" {
		response.Write(w, http.StatusBadRequest, "Place Description cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	res, err := d.status.CheckIn(ctx, body)
	if err != nil {
		e, ok := err.(*myerror.Error)
		if !ok {
			response.Write(w, http.StatusInternalServerError, "Our server encounter a problem.", nil, "BAD-ERROR")
			return
		}
		response.Write(w, http.StatusBadRequest, e.Error(), nil, e.ErrorCode)
		return
	}

	response.Write(w, http.StatusOK, "Checkin Successful", res, "")
}
