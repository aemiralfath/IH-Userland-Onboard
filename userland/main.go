package main

import (
	"os"
	"time"

	"github.com/aemiralfath/IH-Userland-Onboard/userland/api"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/crypto"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/email"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/jwt"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/api/kafka"
	"github.com/aemiralfath/IH-Userland-Onboard/userland/datastore"
	"github.com/rs/zerolog/log"
)

func main() {

	serverCfg := api.ServerConfig{
		Host:            os.Getenv("SERVER_HOST"),
		Port:            os.Getenv("SERVER_PORT"),
		ReadTimeout:     500 * time.Millisecond,
		WriteTimeout:    500 * time.Millisecond,
		ShutdownTimeout: 10 * time.Second,
	}

	emailCfg := email.EmailConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}

	jwtCfg := jwt.JWTConfig{
		Alg:       os.Getenv("JWT_ALG"),
		SignKey:   os.Getenv("JWT_SIGN"),
		VerifyKey: nil,
	}

	kafkaCfg := kafka.KafkaConfig{
		Host:   os.Getenv("KAFKA_HOST"),
		Port:   os.Getenv("KAFKA_PORT"),
		Group:  os.Getenv("KAFKA_GROUP"),
		Offset: os.Getenv("KAFKA_OFFSET"),
	}

	postgresCfg := datastore.PostgresConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
	}

	redisCfg := datastore.RedisConfig{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	log.Info().Msg("get connection to postgres")
	postgresDB, err := datastore.NewPG(postgresCfg)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to postgres")
		return
	}

	log.Info().Msg("get connection to redis")
	redisDB, err := datastore.NewRedis(redisCfg)
	if err != nil {
		log.Error().Err(err).Stack().Msg("failed to connect to redis")
		return
	}

	serverDataSource := &api.DataSource{
		PostgresDB: postgresDB,
		RedisDB:    redisDB,
	}

	serverHelperSouce := &api.HelperSource{
		Jwtauth: jwt.New(jwtCfg),
		Email:   email.NewEmail(emailCfg),
		Crypto:  crypto.NewAppCrypto(),
		Kafka:   kafka.NewKafka(kafkaCfg),
	}

	log.Info().Msg("starting api server")
	srv := api.NewServer(serverCfg, serverHelperSouce, serverDataSource)
	srv.Start()
}
