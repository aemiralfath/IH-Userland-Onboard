package me_test

import (
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/api/jwt/mock"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestUpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := &me.UpdateProfileRequest{
		Fullname: "Awesome User",
		Location: "Jakarta, Indonesia",
		Bio:      "my short bio",
		Web:      "https://example.com",
	}

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("Authorization", "{{TOKEN_TYPE}} {{ACCESS_TOKEN}}")

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(request)

	claims := make(map[string]interface{})
	claims["userID"] = float64(0)

	jwtMock := mock_jwt.NewMockJWT(ctrl)
	jwtMock.EXPECT().
		FromContext(gomock.Any()).
		Times(1).
		Return(nil, claims, nil)

	userId := claims["userID"].(float64)
	profileStore := mock_datastore.NewMockProfileStore(ctrl)
	profileStore.EXPECT().
		UpdateProfile(gomock.Any(), me.ParseUpdateRequestProfile(request), userId).
		Times(1).
		Return(nil)

	handler := me.UpdateProfile(jwtMock, profileStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
