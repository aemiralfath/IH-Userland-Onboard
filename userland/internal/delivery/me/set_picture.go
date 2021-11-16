package me

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/myerror"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
	"github.com/rs/zerolog/log"
)

func (d *DeliveryMe) SetPicture(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.Write(w, http.StatusBadRequest, "File Error.", nil, "BAD-ERROR")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.Write(w, http.StatusBadRequest, "File Error.", nil, "BAD-ERROR")
		return
	}
	defer file.Close()

	userId := claims["userId"]
	fileName := fmt.Sprintf("%s/%s-%s", os.Getenv("PROFILE_PATH"), userId.(string), handler.Filename)
	localFile, err := os.Create(fileName)
	if err != nil {
		log.Error().Err(err).Stack().Msg(err.Error())
		response.Write(w, http.StatusBadRequest, "File Error.", nil, "BAD-ERROR")
		return
	}
	defer localFile.Close()

	if _, err := io.Copy(localFile, file); err != nil {
		log.Error().Err(err).Stack().Msg(err.Error())
		response.Write(w, http.StatusInternalServerError, "Our server encounter a problem.", nil, "BAD-ERROR")
		return
	}

	res, err := d.me.SetPicture(ctx, userId.(string), fileName)
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
