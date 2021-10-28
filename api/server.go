package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/session"
	"github.com/aemiralfath/IH-Userland-Onboard/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/postgres"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/redisdb"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type ServerConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DataSource struct {
	PostgresDB *sql.DB
	RedisDB    *redis.Client
}

type HelperSource struct {
	Jwtauth *jwt.JWTAuth
	Email   *email.Email
}

type serverStores struct {
	userStore     datastore.UserStore
	profileStore  datastore.ProfileStore
	passwordStore datastore.PasswordStore
	sessionStore  datastore.SessionStore
	clientStore   datastore.ClientStore
	otpStore      datastore.OTPStore
}

type Server struct {
	Config     ServerConfig
	DataSource *DataSource
	helper     *HelperSource
	stores     *serverStores
}

func NewServer(config ServerConfig, helper *HelperSource, dataSource *DataSource) *Server {
	return &Server{
		Config:     config,
		DataSource: dataSource,
		helper:     helper,
	}
}

func (s *Server) initStores() error {
	userStore := postgres.NewUserStore(s.DataSource.PostgresDB)
	profileStore := postgres.NewProfileStore(s.DataSource.PostgresDB)
	passwordStore := postgres.NewPasswordStore(s.DataSource.PostgresDB)
	sessionStore := postgres.NewSessionStore(s.DataSource.PostgresDB)
	clientStore := postgres.NewClientStore(s.DataSource.PostgresDB)
	otpStore := redisdb.NewOTPStore(s.DataSource.RedisDB)

	s.stores = &serverStores{
		userStore:     userStore,
		profileStore:  profileStore,
		passwordStore: passwordStore,
		sessionStore:  sessionStore,
		clientStore:   clientStore,
		otpStore:      otpStore,
	}
	return nil
}

func (s *Server) createHandlers() http.Handler {
	// TODO pprof and healthcheck
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(jwt.Verifier(s.helper.Jwtauth))
		r.Use(jwt.Authenticator)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", me.GetProfile(*s.helper.Jwtauth, s.stores.profileStore))
			r.Post("/", me.UpdateProfile(*s.helper.Jwtauth, s.stores.profileStore))

			r.Get("/email", me.GetEmail(*s.helper.Jwtauth, s.stores.userStore))
			r.Post("/email", me.ChangeEmail(*s.helper.Jwtauth, *s.helper.Email, s.stores.userStore, s.stores.otpStore))

			r.Post("/picture", me.SetPicture(*s.helper.Jwtauth, s.stores.profileStore))
			r.Delete("/picture", me.DeletePicture(*s.helper.Jwtauth, s.stores.profileStore))

			r.Post("/password", me.ChangePassword(*s.helper.Jwtauth, s.stores.userStore, s.stores.passwordStore))
			r.Post("/delete", me.DeleteAccount(*s.helper.Jwtauth, s.stores.userStore))

			r.Route("/session", func(r chi.Router) {
				// r.Get("/", session.GetListSession())
				r.Get("/refresh_token", session.GetRefreshToken(*s.helper.Jwtauth))
				r.Get("/access_token", session.GetAccessToken(*s.helper.Jwtauth))
			})
		})

	})

	r.Group(func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hi"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(*s.helper.Email, s.stores.userStore, s.stores.profileStore, s.stores.passwordStore, s.stores.otpStore))
			r.Post("/verification", auth.Verification(*s.helper.Email, s.stores.otpStore, s.stores.userStore))
			r.Post("/login", auth.Login(*s.helper.Jwtauth, s.stores.userStore, s.stores.sessionStore, s.stores.clientStore))

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", auth.ForgotPassword(*s.helper.Email, s.stores.userStore, s.stores.otpStore))
				r.Post("/reset", auth.ResetPassword(s.stores.userStore, s.stores.passwordStore, s.stores.otpStore))
			})
		})
	})

	return r
}

func (s *Server) Start() {
	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(osSigChan)
		os.Exit(0)
	}()

	_ = s.initStores()

	r := s.createHandlers()
	address := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	srv := &http.Server{
		Addr:         address,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
		Handler:      r,
	}

	shutdownCtx := context.Background()
	if s.Config.ShutdownTimeout > 0 {
		var cancelShutdownTimeout context.CancelFunc
		shutdownCtx, cancelShutdownTimeout = context.WithTimeout(shutdownCtx, s.Config.ShutdownTimeout)
		defer cancelShutdownTimeout()
	}

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Stack().Msg("cannot start server")
	}

	log.Info().Msg(fmt.Sprintf("serving %s\n", address))
	go func(srv *http.Server) {
		<-osSigChan
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			panic("failed to shutdown gracefully")
		}
	}(srv)
}
