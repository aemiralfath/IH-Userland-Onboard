package auth_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/api/crypto"
	mock_crypto "github.com/aemiralfath/IH-Userland-Onboard/api/crypto/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	mock_jwt "github.com/aemiralfath/IH-Userland-Onboard/api/jwt/mock"
	mock_kafka "github.com/aemiralfath/IH-Userland-Onboard/api/kafka/mock"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
	"github.com/golang/mock/gomock"
)

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	request := auth.LoginRequest{
		Email:    "account@example.com",
		Password: "-/P4s5w0Rd_!?-WhaT#v#r",
	}

	headers := http.Header{}
	headers.Add("content-type", "application/json")
	headers.Add("X-Api-ClientID", "postman")

	w := &helper.ResponseWriter{}
	r := &http.Request{
		Header: headers,
	}

	r.Body = helper.RequestBody(request)

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
		Return(true)

	jwtToken := &jwt.Token{
		Value:     "{{ACCESS_TOKEN}}",
		Type:      "{{TOKEN_TYPE}}",
		ExpiredAt: time.Now().Add(time.Duration(jwt.AccessTokenExpiration) * time.Minute),
	}
	jwtMock := mock_jwt.NewMockJWT(ctrl)
	jwtMock.EXPECT().
		CreateToken(gomock.Eq(user.ID), gomock.Eq(user.Email), jwt.AccessTokenExpiration).
		Times(1).
		Return(jwtToken, "jti", nil)

	client := &datastore.Client{
		ID:   0,
		Name: r.Header.Get("X-API-ClientID"),
	}

	clientStore := mock_datastore.NewMockClientStore(ctrl)
	clientStore.EXPECT().
		GetClientByName(gomock.Any(), gomock.Eq(client.Name)).
		Times(1).
		Return(client, nil)

	session := &datastore.Session{
		JTI:       "jti",
		UserId:    user.ID,
		IsCurrent: true,
	}

	sessionStore := mock_datastore.NewMockSessionStore(ctrl)
	sessionStore.EXPECT().
		AddNewSession(gomock.Any(), gomock.Eq(session), client.ID).
		Times(1).
		Return(nil)

	kafkaMock := mock_kafka.NewMockKafka(ctrl)

	handler := auth.Login(jwtMock, cryptoMock, kafkaMock, userStore, sessionStore, clientStore)
	handler(w, r)

	result := w.GetBodyJSON()

	if result["require_tfa"] != false {
		t.Errorf("handler not success")
	}

	if len(result) != 2 {
		t.Errorf("handler not success")
	}

}
