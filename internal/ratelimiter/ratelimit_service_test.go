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
func Test_rateLimitService_Allow_DenyWaitAllow(t *testing.T) {
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
	// wait the period window to request again with success
	time.Sleep(time.Second)
	allowResult := ratelimiterService.Allow(context.Background(), ipTest)

	//Then

	assert.Equal(t, true, allowResult, ipTest)
}

func Test_rateLimitUsecase_DenyToken(t *testing.T) {
	//Given

	redisCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisRepository := database.NewRedisRepository(redisCache)
	ratelimiterUsecase := NewRateLimiterUsecase(redisRepository)

	//When
	token := "127.0.0.1"
	input := AllowRateLimitInputDTO{
		Token: token,
		TokenLimit: map[string]int{
			token: 3,
		},
		TokenDuration: time.Second,
	}

	for i := 0; i < 3; i++ {
		ratelimiterUsecase.Execute(context.Background(), input)
	}
	allowResult := ratelimiterUsecase.Execute(context.Background(), input)

	//Then
	assert.Equal(t, false, allowResult.Allow, token)
}
