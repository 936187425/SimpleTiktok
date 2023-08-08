package model

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/go-redis/redis"
)

var RedisHandle *redis.Client

func init() {
	addr := fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port)
	RedisHandle = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	pong, err := RedisHandle.Ping().Result()
	if err == redis.Nil {
		fmt.Printf("Redis 异常！")
	} else if err != nil {
		fmt.Printf("Redis 连接失败！")
	} else {
		fmt.Printf(pong)
	}
}
