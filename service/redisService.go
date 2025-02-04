package service

import (
	"RCSP/global"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisService struct {
	client *redis.Client
	ctx    context.Context
}

// 初始化 Redis 客戶端
func NewRedisService(addr string) *RedisService {
	rdb := global.GvaRedis

	return &RedisService{
		client: rdb,
		ctx:    context.Background(),
	}
}

// 設置值和過期時間
// 在 RedisService 結構體中添加 Set 方法
func (r *RedisService) Set(key string, value string, expire int64) error {
	// 使用超時上下文進行設定
	//timeoutCtx, cancelFunc := context.WithTimeout(r.ctx, time.Second*10)
	timeoutCtx, cancelFunc := context.WithTimeout(r.ctx, time.Second*10)
	defer cancelFunc()

	// 設定值到 Redis
	err := r.client.Set(timeoutCtx, key, value, time.Duration(expire)*time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// 獲取值
func (r *RedisService) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// 刪除值
func (r *RedisService) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}
