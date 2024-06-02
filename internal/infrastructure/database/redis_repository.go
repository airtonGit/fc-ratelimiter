package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string, expiration time.Duration) (int64, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) Exists(ctx context.Context, key string) (bool, error) {

	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		log.Fatalf("Failed to check if key exists: %v", err)
	}
	if exists == 1 {
		return true, nil
	}

	return false, nil
}

func (r *RedisRepository) Incr(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	//newValue, err := r.client.Incr(ctx, key).Result()
	//if err != nil {
	//	log.Fatalf("failed to increment key: %v", err)
	//}
	//
	//return newValue, nil

	// Use Lua script to ensure atomic operations in Redis
	script := redis.NewScript(`
            local count = redis.call("INCR", KEYS[1])
            if count == 1 then
                redis.call("EXPIRE", KEYS[1], ARGV[1])
            end
            return count
        `)

	count, err := script.Run(ctx, r.client, []string{key}, expiration.Seconds()).Result()
	if err != nil {
		return 0, fmt.Errorf("internal Server Error %v", http.StatusInternalServerError)
	}

	return count.(int64), nil
}

func (r *RedisRepository) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}
	return nil
}

//func (r *RedisRepository) BeginTransaction(ctx context.Context) (Cache, error) {
//	r.client.
//}

//func (r *RedisRepository) SetNx(ctx context.Context, key string, value any, expiration time.Duration) error {
//	err := r.client.Set(ctx, key, value, expiration).Err()
//	if err != nil {
//		return fmt.Errorf("failed to set key: %v", err)
//	}
//
//	return nil
//}

// try to setnx, will fail if already exist
