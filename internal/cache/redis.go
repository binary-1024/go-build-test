package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis客户端封装
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		// 如果解析失败，使用默认配置
		opt = &redis.Options{
			Addr: "localhost:6379",
		}
	}

	rdb := redis.NewClient(opt)

	return &RedisClient{
		client: rdb,
	}
}

// Set 设置缓存
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, expiration).Err()
}

// Get 获取缓存
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), dest)
}

// Delete 删除缓存
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (r *RedisClient) Exists(ctx context.Context, key string) bool {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

// Close 关闭连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}
