package ratelimiter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
)

type RateLimitService interface {
	Allow(ctx context.Context, ipToken string) bool
}

type rateLimitService struct {
	cache                 database.Cache
	RequestsByPeriodLimit int
	PeriodDuration        time.Duration
}

func NewRateLimitService(RequestsBySecondLimit int, periodDuration time.Duration, cache database.Cache) RateLimitService {
	return &rateLimitService{
		cache:                 cache,
		RequestsByPeriodLimit: RequestsBySecondLimit,
		PeriodDuration:        periodDuration,
	}
}

func (s *rateLimitService) Allow(ctx context.Context, ipOrToken string) bool {

	if ipOrToken == "" {
		log.Fatalf("Failed to set key: %v", fmt.Errorf(ipOrToken+" is empty"))
	}

	newValue, err := s.cache.Incr(ctx, ipOrToken, s.PeriodDuration)
	if err != nil {
		log.Fatalf("Failed to increment key: %v", err)
	}

	if newValue > int64(s.RequestsByPeriodLimit) {
		return false
	}

	return true
}
