package auth_test

import (
	"net/http"
	"testing"

	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto/mock"
	mock_email "github.com/aemiralfath/IH-Userland-Onboard/userland/api/email/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestForgotPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := auth.ForgotPasswordRequest{
		Email: "account@example.com",
	}

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(request)

	user := &datastore.User{
		ID:    0,
		Email: "account@example.com",
	}

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		GetUserByEmail(gomock.Any(), gomock.Eq(request.Email)).
		Times(1).
		Return(user, nil)

	otpCode := "123456"

	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		GenerateOTP(6).
		Times(1).
		Return(otpCode, nil)

	otpStore := mock_datastore.NewMockOTPStore(ctrl)
	otpStore.EXPECT().
		SetOTP(gomock.Any(), "password", otpCode, request.Email).
		Times(1).
		Return(nil)

	emailMock := mock_email.NewMockEmail(ctrl)

	handler := auth.ForgotPassword(emailMock, cryptoMock, userStore, otpStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
