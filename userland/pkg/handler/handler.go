package handler

import (
	"net/http"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	response.Write(w, http.StatusNotFound, "Route does not exist, please check again your route path.", nil, "Handler-404")
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	response.Write(w, http.StatusMethodNotAllowed, "Method is not allowed, please check again your route method.", nil, "Handler-405")
}
