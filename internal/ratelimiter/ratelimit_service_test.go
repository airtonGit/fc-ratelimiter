package ratelimiter

import (
	"context"
	"testing"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func Test_rateLimitService_Allow_Deny(t *testing.T) {
	//Given

	redisCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisRepository := database.NewRedisRepository(redisCache)
	ratelimiterService := NewRateLimitService(3, time.Second, redisRepository)

	//When
	ipTest := "127.0.0.1"

	for i := 0; i < 3; i++ {
		ratelimiterService.Allow(context.Background(), ipTest)
	}
	allowResult := ratelimiterService.Allow(context.Background(), ipTest)

	//Then

	assert.Equal(t, false, allowResult, ipTest)
}

// Test that hits the limit, waits and try again
func Test_rateLimitService_Allow_RequestWaitAllow(t *testing.T) {
	//Given

	redisCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisRepository := database.NewRedisRepository(redisCache)
	ratelimiterService := NewRateLimitService(3, time.Second, redisRepository)

	//When
	ipTest := "127.0.0.1"

	for i := 0; i < 2; i++ {
		ratelimiterService.Allow(context.Background(), ipTest)
	}

	// wait the period window to request again with success
	time.Sleep(time.Second)
	allowResult := ratelimiterService.Allow(context.Background(), ipTest)

	//Then
	assert.Equal(t, true, allowResult, ipTest)
}
