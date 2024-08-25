package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

var rdb *redis.Client

func init() {
	addr := "127.0.0.1"
	port := "6379"
	password := ""
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr + ":" + port,
		Password: password,
		DB:       0,
		PoolSize: 100,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		logrus.Error("redis ping err:", err)
		panic("redis连接失败：" + err.Error())
	}
}

func GetRedisDb() *redis.Client {
	return rdb
}
