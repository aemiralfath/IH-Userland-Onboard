package status

import (
	"net/http"
	"strings"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
	"github.com/go-chi/chi"
)

func (d *DeliveryStatus) Report(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	placeId := chi.URLParam(r, "placeId")

	if strings.TrimSpace(placeId) == "" {
		response.Write(w, http.StatusBadRequest, "Place id cannot be empty.", nil, "AUTH-DLV-02")
		return
	}

	res, err := d.status.Report(ctx, placeId)
	if err != nil {
		e, ok := err.(*myerror.Error)
		if !ok {
			response.Write(w, http.StatusInternalServerError, "Our server encounter a problem.", nil, "BAD-ERROR")
			return
		}
		response.Write(w, http.StatusBadRequest, e.Error(), nil, e.ErrorCode)
		return
	}

	response.Write(w, http.StatusOK, "Get Report Successful", res, "")
}
