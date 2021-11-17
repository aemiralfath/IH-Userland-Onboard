package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestResetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := auth.ResetPasswordRequest{
		Token:           "RESET_PASSWORD_TOKEN",
		Password:        "-/P4s5w0Rd_!?-N#w",
		PasswordConfirm: "-/P4s5w0Rd_!?-N#w",
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

	otpStore := mock_datastore.NewMockOTPStore(ctrl)
	otpStore.EXPECT().
		GetOTP(gomock.Any(), "password", request.Token).
		Times(1).
		Return(user.Email, nil)

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
		Times(1).
		Return(user, nil)

	lastThreePassword := []string{}
	for i := 0; i < 2; i++ {
		userPassword, _ := crypto.NewAppCrypto().HashPassword(fmt.Sprintf("-/P4s5w0Rd_!?-WhaT#v#r%d", i))
		lastThreePassword = append(lastThreePassword, userPassword)
	}

	passwordStore := mock_datastore.NewMockPasswordStore(ctrl)
	passwordStore.EXPECT().
		GetLastThreePassword(gomock.Any(), gomock.Eq(user.ID)).
		Times(1).
		Return(lastThreePassword, nil)

	mockCrypto := mock_crypto.NewMockCrypto(ctrl)
	for _, e := range lastThreePassword {
		mockCrypto.EXPECT().
			ConfirmPassword(e, request.Password).
			Return(false)
	}

	hashPassword, _ := crypto.NewAppCrypto().HashPassword(request.Password)
	mockCrypto.EXPECT().
		HashPassword(request.Password).
		Times(1).
		Return(hashPassword, nil)

	userStore.EXPECT().
		ChangePassword(gomock.Any(), user).
		Times(1).
		Return(nil)

	passwordStore.EXPECT().
		AddNewPassword(gomock.Any(), &datastore.Password{Password: hashPassword}, user.ID).
		Times(1).
		Return(nil)

	handler := auth.ResetPassword(mockCrypto, userStore, passwordStore, otpStore)
	handler(w, r)

	result := w.GetBodyJSON()
	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
