package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedisConnection() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	addr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  redisPassword,
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Error().Err(err).Str("addr", addr).Msg("Failed to connect to Redis")
	} else {
		log.Info().Str("addr", addr).Msg("Successfully connected to Redis")
	}
}

func StopRedisConnection() {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing Redis connection")
		} else {
			log.Info().Msg("Success closing connection to Redis")
		}
	}
}
