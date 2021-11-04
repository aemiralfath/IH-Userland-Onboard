package me_test

import (
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/api/jwt/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := http.Header{}
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

	profile := &datastore.Profile{
		ID:        0,
		Fullname:  "Awesome User",
		Location:  "Jakarta, Indonesia",
		Bio:       "my short bio",
		Web:       "https://userland-api.simukti.net",
		Picture:   "https://userlandapp.storage.googleapis.com/us",
		CreatedAt: "2009-11-10T23:00:00Z",
	}

	userId := claims["userID"].(float64)
	profileStore := mock_datastore.NewMockProfileStore(ctrl)
	profileStore.EXPECT().
		GetProfile(gomock.Any(), userId).
		Times(1).
		Return(profile, nil)

	handler := me.GetProfile(jwtMock, profileStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if len(result) != 7 {
		t.Errorf("handler not success")
	}
}
