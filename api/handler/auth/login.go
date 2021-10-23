package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(jwtAuth helper.JWTAuth, userStore datastore.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		req := &loginRequest{}

		if err := render.Bind(r, req); err != nil {
			fmt.Println(render.Render(w, r, handler.BadRequestErrorRenderer(err)))
			return
		}

		usr, err := userStore.GetUser(ctx, parseLoginUser(req))
		if err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		fmt.Printf("%d %s %s\n", usr.ID, usr.Email, usr.Password)

		if err := confirmPassword(usr.Password, req.Password); err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		accessTokenClaims := make(map[string]interface{})
		accessTokenClaims["id"] = usr.ID
		accessTokenClaims["email"] = usr.Email
		helper.SetIssuedNow(accessTokenClaims)
		helper.SetExpiryIn(accessTokenClaims, time.Duration(helper.AccessTokenExpiration))

		_, accessToken, err := jwtAuth.Encode(accessTokenClaims)
		if err != nil {
			fmt.Println(render.Render(w, r, handler.InternalServerErrorRenderer(err)))
			return
		}

		// refreshTokenClaims := make(map[string]interface{})
		// refreshTokenClaims["email"] = usr.Email
		// helper.SetIssuedNow(refreshTokenClaims)
		// helper.SetExpiryIn(refreshTokenClaims, time.Duration(helper.RefreshTokenExpiration))

		// _, refreshToken, err := jwtAuth.Encode(refreshTokenClaims)
		// if err != nil {
		// 	fmt.Println(render.Render(w, r, handler.BadRequestErrorRenderer(err)))
		// 	return
		// }

		handler.CustomRender(w, http.StatusOK, map[string]interface{}{
			"require_tfa": false,
			"access_token": map[string]string{
				"value":      accessToken,
				"type":       "BEARER",
				"expired_at": time.Unix(accessTokenClaims["exp"].(int64), 0).String(),
			},
		})
	}
}

func parseLoginUser(u *loginRequest) *datastore.User {
	return &datastore.User{
		Email:    u.Email,
		Password: u.Password,
	}
}

func confirmPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (login *loginRequest) Bind(r *http.Request) error {
	if strings.TrimSpace(login.Email) == "" {
		return fmt.Errorf("required email")
	}

	if strings.TrimSpace(login.Password) == "" {
		return fmt.Errorf("required password")
	}

	return nil
}

func (*loginRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
