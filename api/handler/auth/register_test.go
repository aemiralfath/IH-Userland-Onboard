package auth_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/api/crypto/mock"
	mock_email "github.com/aemiralfath/IH-Userland-Onboard/api/email/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	request := auth.RegisterRequest{
		Fullname:        "Account Fullname",
		Email:           "account@example.com",
		Password:        "-/P4s5w0Rd_!?-WhaT#v#r",
		PasswordConfirm: "-/P4s5w0Rd_!?-WhaT#v#r",
	}

	user := &datastore.User{
		ID:       0,
		Email:    "account@example.com",
		Password: "-/P4s5w0Rd_!?-WhaT#v#r",
	}

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		CheckUserEmailExist(gomock.Any(), gomock.Eq(user.Email)).
		Times(1).
		Return(nil, sql.ErrNoRows)

	cryptoService := &crypto.AppCrypto{}
	hashPassword, _ := cryptoService.HashPassword(user.Password)

	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		HashPassword(gomock.Eq(user.Password)).
		Times(1).
		Return(hashPassword, nil)

	user.Password = hashPassword

	userStore.EXPECT().
		AddNewUser(gomock.Any(), gomock.Eq(user)).
		Times(1).
		Return(user.ID, nil)

	profileStore := mock_datastore.NewMockProfileStore(ctrl)
	profileStore.EXPECT().
		AddNewProfile(gomock.Any(), &datastore.Profile{Fullname: request.Fullname}, user.ID).
		Times(1).
		Return(nil)

	passwordStore := mock_datastore.NewMockPasswordStore(ctrl)
	passwordStore.EXPECT().
		AddNewPassword(gomock.Any(), &datastore.Password{Password: hashPassword}, user.ID).
		Times(1).
		Return(nil)

	otpCode := "123456"
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

	subject := "Userland Email Verification!"
	msg := fmt.Sprintf("Use this otp for verify your email: %s", otpCode)


	emailMock := mock_email.NewMockEmail(ctrl)
	emailMock.EXPECT().
		SendEmail(request.Email, subject, msg).
		Times(1)

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(request)

	handler := auth.Register(
		emailMock,
		cryptoMock,
		userStore,
		profileStore,
		passwordStore,
		otpStore)

	handler(w, r)

	result := w.GetBodyJSON()

	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
