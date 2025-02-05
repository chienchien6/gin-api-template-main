// Package initialize /**
package initialize

import (
	"RCSP/global"
	"context"

	"github.com/go-redis/redis/v8"
)

// InitRedis 初始化redis客户端
func InitRedis() {
	if !global.GvaConfig.Redis.Enable {
		return
	}
	// 创建
	redisClient := redis.NewClient(&redis.Options{
		Addr:     global.GvaConfig.Redis.Addr,
		Password: global.GvaConfig.Redis.Password,
		DB:       global.GvaConfig.Redis.DefaultDB,
	})
	// 使用超时上下文，验证redis
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), global.GvaConfig.Redis.DialTimeout)
	defer cancelFunc()
	_, err := redisClient.Ping(timeoutCtx).Result() //此方法用於檢查與 Redis 服務器的連接是否正常。它會發送一個 PING 命令到 Redis，並期望收到 PONG 響應。
	if err != nil {
		panic("redis初始化失败! " + err.Error())
	}
	global.GvaRedis = redisClient

}
