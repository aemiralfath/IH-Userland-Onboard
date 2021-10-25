package main

import (
	"log"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/api"
	"github.com/aemiralfath/IH-Userland-Onboard/datastore"
)

func main() {
	// TODO use external config management (toml?)
	serverCfg := api.ServerConfig{
		Host:            "0.0.0.0",
		Port:            "8080",
		ReadTimeout:     500 * time.Millisecond,
		WriteTimeout:    500 * time.Millisecond,
		ShutdownTimeout: 10 * time.Second,
	}

	postgresCfg := datastore.PostgresConfig{
		Host:     "db_userland",
		Port:     5432,
		Username: "admin",
		Password: "admin",
		Database: "userland",
	}

	redisCfg := datastore.RedisConfig{
		Address:  "redis_userland:6379",
		Password: "",
		DB:       0,
	}

	postgresDB, err := datastore.NewPG(postgresCfg)
	if err != nil {
		// TODO proper logging with zlogger
		log.Fatalf("failed to open db conn: %v\n", err)
	}

	redisDB, err := datastore.NewRedis(redisCfg)
	if err != nil {
		// TODO proper logging with zlogger
		log.Fatalf("failed to open redis conn: %v\n", err)
	}

	serverDataSource := &api.DataSource{
		PostgresDB: postgresDB,
		RedisDB:    redisDB,
	}

	srv := api.NewServer(serverCfg, serverDataSource)
	srv.Start()
}
