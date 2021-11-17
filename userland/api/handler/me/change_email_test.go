package me_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto/mock"
	mock_email "github.com/aemiralfath/IH-Userland-Onboard/userland/api/email/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt/mock"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestChangeEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := me.ChangeEmailRequest{
		Email: "newemail@example.com",
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

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		CheckUserEmailExist(gomock.Any(), request.Email).
		Times(1).
		Return(nil, sql.ErrNoRows)

	otpCode := "123456"
	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		GenerateOTP(6).
		Return(otpCode, nil)

	userId := claims["userID"].(float64)
	otpValue := fmt.Sprintf("%f-%s", userId, request.Email)
	otpMock := mock_datastore.NewMockOTPStore(ctrl)
	otpMock.EXPECT().
		SetOTP(gomock.Any(), "user", otpCode, otpValue).
		Times(1).
		Return(nil)

	emailMock := mock_email.NewMockEmail(ctrl)

	handler := me.ChangeEmail(jwtMock, cryptoMock, emailMock, userStore, otpMock)
	handler(w, r)

	result := w.GetBodyJSON()

	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
