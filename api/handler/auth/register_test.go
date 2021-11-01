package auth_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/crypto"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	email := mockEmail()
	user := mockUser()

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		CheckUserEmailExist(gomock.Any(), gomock.Eq(user.Email)).
		Times(1).
		Return(nil, sql.ErrNoRows)

	cryptoService := &crypto.AppCrypto{}
	hashPassword, _ := cryptoService.HashPassword(user.Password)

	fmt.Println(hashPassword)
	cryptoMock := mock_datastore.NewMockCrypto(ctrl)
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
		AddNewProfile(gomock.Any(), &datastore.Profile{Fullname: "Test"}, user.ID).
		Times(1).
		Return(nil)

	passwordStore := mock_datastore.NewMockPasswordStore(ctrl)
	passwordStore.EXPECT().
		AddNewPassword(gomock.Any(), &datastore.Password{Password: "Test1234"}, user.ID).
		Times(1).
		Return(nil)

	otpStore := mock_datastore.NewMockOTPStore(ctrl)
	otpStore.EXPECT().
		SetOTP(gomock.Any(), "user", "123456", user.Email).
		Times(1).
		Return(nil)

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(map[string]string{
		"fullname":         "Test",
		"email":            "test@gmail.com",
		"password":         "Test1234",
		"password_confirm": "Test1234",
	})

	handler := auth.Register(*email, cryptoService, userStore, profileStore, passwordStore, otpStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if len(result) != 1 {
		t.Errorf("Item was not added to the datastore")
	}
}

func mockUser() *datastore.User {
	return &datastore.User{
		ID:       0,
		Email:    "test@gmail.com",
		Password: "Test1234",
	}
}

func mockEmail() *email.Email {
	return email.NewEmail(email.EmailConfig{
		Host:     "smtp-relay.sendinblue.com",
		Port:     "587",
		From:     "35IKwbUMAygQH6Ch",
		Password: "criptdestroyer@gmail.com",
	})
}
