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

	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/auth"
	"github.com/aemiralfath/IH-Userland-Onboard/api/handler/me"
	"github.com/aemiralfath/IH-Userland-Onboard/api/helper"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/postgres"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore/redisdb"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis/v8"
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

type stores struct {
	userStore     datastore.UserStore
	profileStore  datastore.ProfileStore
	passwordStore datastore.PasswordStore
	otpStore      datastore.TokenStore
}

type jwtConfig struct {
	tokenAuth *helper.JWTAuth
}

type Server struct {
	Config     ServerConfig
	DataSource *DataSource
	stores     *stores
	jwt        *jwtConfig
}

func NewServer(config ServerConfig, dataSource *DataSource) *Server {
	return &Server{
		Config:     config,
		DataSource: dataSource,
	}
}

func (s *Server) initStores() error {
	userStore := postgres.NewUserStore(s.DataSource.PostgresDB)
	profileStore := postgres.NewProfileStore(s.DataSource.PostgresDB)
	passwordStore := postgres.NewPasswordStore(s.DataSource.PostgresDB)
	otpStore := redisdb.NewTokenStore(s.DataSource.RedisDB)
	s.stores = &stores{
		userStore:     userStore,
		profileStore:  profileStore,
		passwordStore: passwordStore,
		otpStore:      otpStore,
	}
	return nil
}

func (s *Server) initJwt() error {
	// user .env for secret
	tokenAuth := helper.New("HS256", []byte("secretSign"), nil)
	s.jwt = &jwtConfig{
		tokenAuth: tokenAuth,
	}
	return nil
}

func (s *Server) createHandlers() http.Handler {
	// TODO pprof and healthcheck
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(helper.Verifier(s.jwt.tokenAuth))
		r.Use(helper.Authenticator)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", me.GetProfile(*s.jwt.tokenAuth, s.stores.profileStore))
			r.Post("/", me.UpdateProfile(*s.jwt.tokenAuth, s.stores.profileStore))
		})
		
	})

	r.Group(func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("Hi"))
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", auth.Register(s.stores.userStore, s.stores.profileStore, s.stores.passwordStore, s.stores.otpStore))
			r.Post("/verification", auth.Verification(s.stores.otpStore))
			r.Post("/login", auth.Login(*s.jwt.tokenAuth, s.stores.userStore))

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", auth.ForgotPassword(s.stores.userStore, s.stores.otpStore))
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
	_ = s.initJwt()

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
		// TODO replace with zlogger fatal
		panic("cannot start server")
	}
	// TODO with proper logging with zlogger
	fmt.Printf("serving %s\n", address)

	go func(srv *http.Server) {
		<-osSigChan
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			panic("failed to shutdown gracefully")
		}
	}(srv)
}
