package auth_test

import (
	"net/http"
	"testing"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	mock_datastore "github.com/aemiralfath/IH-Userland-Onboard/datastore/mock"
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

	// w := &helper.ResponseWriter{}
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

}
