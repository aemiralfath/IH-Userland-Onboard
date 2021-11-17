package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/app"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/postgres"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/internal/model/redis"
	"github.com/rs/zerolog/log"
)

func main() {
	serverCfg := app.AppConfig{
		Host:            os.Getenv("SERVER_HOST"),
		Port:            os.Getenv("SERVER_PORT"),
		ReadTimeout:     500 * time.Millisecond,
		WriteTimeout:    500 * time.Millisecond,
		ShutdownTimeout: 10 * time.Second,
	}

	postgresCfg := postgres.PostgresConfig{
		Host:     os.Getenv("POSTGRES_ADDR"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	redisCfg := redis.RedisConfig{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	app := app.New(serverCfg, postgresCfg, redisCfg)
	log.Info().Msg(fmt.Sprintf("Starting api server in %s:%s", serverCfg.Host, serverCfg.Port))
	app.StartServer()
}
