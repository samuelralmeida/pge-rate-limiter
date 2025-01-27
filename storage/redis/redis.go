package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage() *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
		Protocol: 2,
	})
	return &RedisStorage{client: client}
}

// padrão recomendado pelo próprio redis: https://redis.io/docs/latest/commands/incr/
func (r *RedisStorage) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	ts := time.Now().Unix()
	key = fmt.Sprintf("%s:%d", key, ts)

	pipe := r.client.Pipeline()

	result := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, fmt.Errorf("error to exec redis pipeline: %w", err)
	}

	fmt.Println(key)

	return int(result.Val()), nil
}

func (r *RedisStorage) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisStorage) Exists(ctx context.Context, key string) (int, error) {
	resp, err := r.client.Exists(ctx, key).Result()
	return int(resp), err
}
