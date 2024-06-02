package ratelimiter

import (
	"context"
	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_rateLimitUsecase_DenyToken(t *testing.T) {
	//Given
	redisCache := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	redisRepository := database.NewRedisRepository(redisCache)
	ratelimiterUsecase := NewRateLimiterUsecase(redisRepository, nil)

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
