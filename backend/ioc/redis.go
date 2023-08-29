package ioc

import (
	"github.com/JhonWong/webook/backend/config"
	"github.com/redis/go-redis/v9"
)

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
