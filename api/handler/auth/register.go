package auth

import (
	"encoding/json"
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

func Register(authStore datastore.AuthStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_ = authStore.GetAuth(ctx)
		success := struct {
			Success bool `json:"success"`
		}{Success: true}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(success)
	}
}