package config

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS.ADDR"),
		Password: viper.GetString("REDIS.PASSWORD"), // no password set
		DB:       viper.GetInt("REDIS.DB"),          // use default DB
	})
}
