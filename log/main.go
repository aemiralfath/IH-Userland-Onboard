package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

func main() {

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_ADDR"),
		os.Getenv("POSTGRES_DB"),
	)

	connectionCfg, err := pgx.ParseConfig(connString)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse config")
	}

	connStr := stdlib.RegisterConnConfig(connectionCfg)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Error().Err(err).Msg("failed to create pgx")
	}

	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("failed to create ping db")
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka-1:19092",
		"group.id":          "myLogin",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to create kafka consumer")
		panic("kafka consumer failed")
	}

	err = c.SubscribeTopics([]string{"login-succeed"}, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to load topics")
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var m map[string]interface{}
			json.Unmarshal(msg.Value, &m)

			id, err := uuid.NewRandom()
			if err != nil {
				log.Error().Err(err).Msg("failed to create id")
				return
			}

			log.Info().Msg(fmt.Sprintf("Adding Data %s, %s, %s", m["userid"], m["remote-ip"], m["username"]))
			query := `INSERT INTO "audit_logs" (id, user_id, remote_ip, username) VALUES ($1, $2, $3, $4)`
			_, err = db.Exec(query, id, m["userid"], m["remote-ip"], m["username"])
			if err != nil {
				log.Error().Err(err).Msg("failed to add message")
				return
			}
			log.Info().Msg("Add log success")
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
