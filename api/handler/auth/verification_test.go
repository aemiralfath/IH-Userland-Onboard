package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/api/crypto/mock"
	mock_email "github.com/aemiralfath/IH-Userland-Onboard/api/email/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestVerification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	request := auth.VerificationRequest{
		Type:      "email.verify",
		Recipient: "account@example.com",
	}

	user := &datastore.User{
		ID:    0,
		Email: "account@example.com",
	}

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		CheckUserEmailExist(gomock.Any(), gomock.Eq(request.Recipient)).
		Times(1).
		Return(user, nil)

	otpCode := "123456"
	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		GenerateOTP(gomock.Eq(6)).
		Times(1).
		Return(otpCode, nil)

	value := fmt.Sprintf("%f-%s", float64(user.ID), user.Email)
	otpStore := mock_datastore.NewMockOTPStore(ctrl)
	otpStore.EXPECT().
		SetOTP(gomock.Any(), "user", gomock.Eq(otpCode), gomock.Eq(value)).
		Times(1).
		Return(nil)

	emailMock := mock_email.NewMockEmail(ctrl)

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(request)

	handler := auth.Verification(
		emailMock,
		cryptoMock,
		otpStore,
		userStore)

	handler(w, r)

	result := w.GetBodyJSON()
	fmt.Println()

	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
