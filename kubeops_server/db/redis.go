package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"kubeops/config"
	"kubeops/utils"
)

var Rdb *redis.Client

func CloseRedis() error {
	return Rdb.Close()
}

func InitRedis() {
	ctx := context.Background()
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost,
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
	})
	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		utils.Logger.Error("Failed to initialize Redis", err.Error())
		panic(err)
	}
	utils.Logger.Info("redis initialization succeeded", pong)
}
