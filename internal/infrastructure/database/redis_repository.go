package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Exists(ctx context.Context, key string) (bool, error)
	Incr(ctx context.Context, key string) (int64, error)
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

func (r *RedisRepository) Incr(ctx context.Context, key string) (int64, error) {
	newValue, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		log.Fatalf("failed to increment key: %v", err)
	}

	//r.client.Eval(ctx, `
	//	FUNCTION LIMIT_API_CALL(ip)
	//		ts = CURRENT_UNIX_TIME()
	//		keyname = ip+":"+ts
	//		MULTI
	//			INCR(keyname)
	//			EXPIRE(keyname,10)
	//		EXEC
	//		current = RESPONSE_OF_INCR_WITHIN_MULTI
	//		IF current > 10 THEN
	//			ERROR "too many requests per second"
	//		ELSE
	//			PERFORM_API_CALL()
	//		END
	//`, []string{key})

	return newValue, nil
}

func (r *RedisRepository) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %v", err)
	}

	// r.client.Set()

	return nil
}
