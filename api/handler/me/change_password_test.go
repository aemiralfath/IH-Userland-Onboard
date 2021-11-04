package me_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/api/crypto/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/api/jwt/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := me.ChangePasswordRequest{
		PasswordCurrent: "-/P4s5w0Rd_!?-N#w",
		Password:        "-/P4s5w0Rd_!?-N#w22",
		PasswordConfirm: "-/P4s5w0Rd_!?-N#w22",
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
	claims["email"] = "account@example.com"

	cryptoService := &crypto.AppCrypto{}
	hashPassword, _ := cryptoService.HashPassword(request.PasswordCurrent)

	user := &datastore.User{
		ID:       0,
		Email:    "account@example.com",
		Password: hashPassword,
	}

	jwtMock := mock_jwt.NewMockJWT(ctrl)
	jwtMock.EXPECT().
		FromContext(gomock.Any()).
		Times(1).
		Return(nil, claims, nil)

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		GetUserByEmail(gomock.Any(), claims["email"].(string)).
		Times(1).
		Return(*user, nil)

	fmt.Println(user.Password)
	fmt.Println(hashPassword)

	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		ConfirmPassword(user.Password, request.PasswordCurrent).
		Times(1).
		Return(true)

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

	for _, e := range lastThreePassword {
		cryptoMock.EXPECT().
			ConfirmPassword(e, request.Password).
			Return(false)
	}

	hashPassword, _ = crypto.NewAppCrypto().HashPassword(request.Password)
	cryptoMock.EXPECT().
		HashPassword(request.Password).
		Times(1).
		Return(hashPassword, nil)

	user.Password = hashPassword
	userStore.EXPECT().
		ChangePassword(gomock.Any(), user).
		Times(1).
		Return(nil)

	passwordStore.EXPECT().
		AddNewPassword(gomock.Any(), &datastore.Password{Password: hashPassword}, user.ID).
		Times(1).
		Return(nil)

	handler := me.ChangePassword(
		jwtMock,
		cryptoMock,
		userStore,
		passwordStore)

	handler(w, r)

	result := w.GetBodyJSON()
	if result["success"] != true {
		t.Errorf("handler not success")
	}
}
