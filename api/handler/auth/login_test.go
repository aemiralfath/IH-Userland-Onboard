package auth_test

import (
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/api/crypto/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	request := auth.LoginRequest{
		Email:    "account@example.com",
		Password: "-/P4s5w0Rd_!?-WhaT#v#r",
	}

	userPassword, _ := crypto.NewAppCrypto().HashPassword("-/P4s5w0Rd_!?-WhaT#v#r")

	user := &datastore.User{
		ID:       0,
		Email:    "account@example.com",
		Password: userPassword,
	}

	userStore := mock_datastore.NewMockUserStore(ctrl)
	userStore.EXPECT().
		GetUserByEmail(gomock.Any(), gomock.Eq(request.Email)).
		Times(1).
		Return(user, nil)

	cryptoMock := mock_crypto.NewMockCrypto(ctrl)
	cryptoMock.EXPECT().
		ConfirmPassword(gomock.Eq(user.Password), gomock.Eq(request.Password)).
		Times(1).
		Return(nil)

}
