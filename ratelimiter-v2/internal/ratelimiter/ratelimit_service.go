package ratelimiter

import (
	"context"
	"log"
	"time"

	"ratelimiter-v2/internal/infra"
)

type RateLimitService interface {
	Allow(ctx context.Context, ip string) bool
}

type rateLimitService struct {
	cache      infra.Cache
	MaxRequest int
	Duration   time.Duration
}

func NewRateLimitService() RateLimitService {
	return &rateLimitService{}
}

func (s *rateLimitService) Allow(ctx context.Context, ipOrToken string) bool {

	// confiro se key já existe (IP ou Token)
	// se nao existe, crio e adiciono ttl
	// se existe incremento, sem mexer no expire
	exists, err := s.cache.Exists(ctx, ipOrToken)
	if err != nil {
		log.Fatalf("Failed to check if key exists: %v", err)
	}
	if exists {
		newValue, err := s.cache.Incr(ctx, ipOrToken)
		if err != nil {
			log.Fatalf("Failed to increment key: %v", err)
		}

		if newValue > int64(s.MaxRequest) {
			return false
		}

	} else {
		// first for this key, in the time window
		err = s.cache.Set(ctx, ipOrToken, 1, s.Duration)
		if err != nil {
			log.Fatalf("Failed to set key: %v", err)
		}
	}

	return true
}
