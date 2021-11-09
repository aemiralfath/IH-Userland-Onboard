package me_test

import (
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt/mock"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestGetEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("Authorization", "{{TOKEN_TYPE}} {{ACCESS_TOKEN}}")

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	claims := make(map[string]interface{})
	claims["userID"] = float64(0)

	jwtMock := mock_jwt.NewMockJWT(ctrl)
	jwtMock.EXPECT().
		FromContext(gomock.Any()).
		Times(1).
		Return(nil, claims, nil)

	email := "user@example.com"

	userId := claims["userID"].(float64)
	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		GetEmailByID(gomock.Any(), userId).
		Times(1).
		Return(email, nil)

	handler := me.GetEmail(jwtMock, userStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if result["email"] != email {
		t.Errorf("handler not success")
	}
}
