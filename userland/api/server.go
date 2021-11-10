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

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/handler/session"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/kafka"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/postgres"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore/redisdb"
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
	Jwtauth jwt.JWT
	Email   email.Email
	Crypto  crypto.Crypto
	Kafka   kafka.Kafka
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
		r.Use(s.helper.Jwtauth.Verifier())
		r.Use(s.helper.Jwtauth.Authenticator)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", me.GetProfile(s.helper.Jwtauth, s.stores.profileStore))
			r.Post("/", me.UpdateProfile(s.helper.Jwtauth, s.stores.profileStore))

			r.Get("/email", me.GetEmail(s.helper.Jwtauth, s.stores.userStore))
			r.Post("/email", me.ChangeEmail(s.helper.Jwtauth, s.helper.Crypto, s.helper.Email, s.stores.userStore, s.stores.otpStore))

			r.Post("/picture", me.SetPicture(s.helper.Jwtauth, s.stores.profileStore))
			r.Delete("/picture", me.DeletePicture(s.helper.Jwtauth, s.stores.profileStore))

			r.Post("/password", me.ChangePassword(s.helper.Jwtauth, s.helper.Crypto, s.stores.userStore, s.stores.passwordStore))
			r.Post("/delete", me.DeleteAccount(s.helper.Jwtauth, s.helper.Crypto, s.stores.userStore))

			r.Route("/session", func(r chi.Router) {
				r.Get("/", session.GetListSession(s.helper.Jwtauth, s.stores.sessionStore))
				r.Delete("/", session.EndCurrentSession(s.helper.Jwtauth, s.stores.sessionStore))
				r.Delete("/other", session.DeleteOtherSession(s.helper.Jwtauth, s.stores.sessionStore))
				r.Get("/refresh_token", session.GetRefreshToken(s.helper.Jwtauth))
				r.Get("/access_token", session.GetAccessToken(s.helper.Jwtauth))
			})
		})

	})

	r.Group(func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hi"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(s.helper.Email, s.helper.Crypto, s.stores.userStore, s.stores.profileStore, s.stores.passwordStore, s.stores.otpStore))
			r.Post("/verification", auth.Verification(s.helper.Email, s.helper.Crypto, s.stores.otpStore, s.stores.userStore))
			r.Post("/login", auth.Login(s.helper.Jwtauth, s.helper.Crypto, s.helper.Kafka, s.stores.userStore, s.stores.sessionStore, s.stores.clientStore))

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", auth.ForgotPassword(s.helper.Email, s.helper.Crypto, s.stores.userStore, s.stores.otpStore))
				r.Post("/reset", auth.ResetPassword(s.helper.Crypto, s.stores.userStore, s.stores.passwordStore, s.stores.otpStore))
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

	// consumer, err := s.helper.Kafka.NewConsumer()
	// if err != nil {
	// 	log.Fatal().Err(err).Stack().Msg("cannot add new consumer")
	// }

	// err = consumer.SubscribeTopics([]string{"login-succeed"}, nil)
	// if err != nil {
	// 	log.Fatal().Err(err).Stack().Msg("cannot subscribe topics")
	// }

	// for {
	// 	msg, err := consumer.ReadMessage(-1)
	// 	if err == nil {
	// 		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
	// 	} else {
	// 		// The client will automatically try to recover from all errors.
	// 		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	// 	}
	// }

	// consumer.Close()
}