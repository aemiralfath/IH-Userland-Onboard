package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/postgres"
	myredis "github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/redis"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/handler"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/pkg/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type AppConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type App struct {
	Config AppConfig
	DB     *sql.DB
	Redis  *redis.Client
}

func New(serverCfg AppConfig, postgresCfg postgres.PostgresConfig, redisCfg myredis.RedisConfig) App {

	log.Info().Msg("get connection to postgres")
	postgre, err := postgres.NewPG(postgresCfg)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to postgres")
	}

	log.Info().Msg("get connection to redis")
	redis, err := myredis.NewRedis(redisCfg)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to redis")
	}

	return App{
		Config: serverCfg,
		DB:     postgre,
		Redis:  redis,
	}
}

func (app *App) createHandlers() http.Handler {
	r := chi.NewRouter()
	d, err := initDelivery(app.DB, app.Redis)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to kafka")
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.NotFound(handler.NotFound)
	r.MethodNotAllowed(handler.MethodNotAllowed)

	r.Group(func(r chi.Router) {
		r.Use(jwt.New().Verifier())
		r.Use(jwt.New().Authenticator)
		r.Route("/me", func(r chi.Router) {
			r.Get("/", d.me.Profile)
			r.Post("/", d.me.UpdateProfile)

			r.Get("/email", d.me.Email)
			r.Post("/email", d.me.ChangeEmail)

			r.Post("/picture", d.me.SetPicture)
			r.Delete("/picture", d.me.DeletePicture)

			r.Post("/password", d.me.ChangePassword)
			r.Post("/delete", d.me.DeleteAccount)

			r.Route("/session", func(r chi.Router) {
				r.Get("/", d.session.ListSession)
				r.Delete("/", d.session.EndCurrentSession)
				r.Delete("/other", d.session.DeleteOtherSession)
				r.Get("/refresh_token", d.session.RefreshToken)
				r.Get("/access_token", d.session.AccessToken)
			})
		})
	})

	r.Group(func(r chi.Router) {
		r.Get("/", func(rw http.ResponseWriter, r *http.Request) {
			response.Write(rw, http.StatusOK, "Hi", nil, "")
		})

		r.Get("/report/{placeId}", d.status.Report)
		r.Post("/checkin", d.status.CheckIn)
		r.Post("/checkout", d.status.CheckOut)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", d.auth.Register)
			r.Post("/verification", d.auth.Verification)
			r.Post("/login", d.auth.Login)

			r.Route("/password", func(r chi.Router) {
				r.Post("/forgot", d.auth.ForgotPassword)
				r.Post("/reset", d.auth.ResetPassword)
			})
		})
	})

	return r
}

func (app *App) StartServer() {
	osSignalChan := make(chan os.Signal, 1)
	signal.Notify(osSignalChan, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(osSignalChan)
		os.Exit(0)
	}()

	r := app.createHandlers()
	address := fmt.Sprintf("%s:%s", app.Config.Host, app.Config.Port)
	server := &http.Server{
		Addr:         address,
		ReadTimeout:  app.Config.ReadTimeout,
		WriteTimeout: app.Config.WriteTimeout,
		Handler:      r,
	}

	shutdownCtx := context.Background()
	if app.Config.ShutdownTimeout > 0 {
		var cancelShutdownTimeout context.CancelFunc
		shutdownCtx, cancelShutdownTimeout = context.WithTimeout(shutdownCtx, app.Config.ShutdownTimeout)
		defer cancelShutdownTimeout()
	}

	err := server.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Stack().Msg("cannot start server")
	}

	log.Info().Msg(fmt.Sprintf("serving %s\n", address))
	go func(srv *http.Server) {
		<-osSignalChan
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			panic("failed to shutdown gracefully")
		}
	}(server)
}
