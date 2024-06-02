package ratelimiter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/go-redsync/redsync/v4"
)

type RateLimitService interface {
	Allow(ctx context.Context, ipToken string) bool
}

type rateLimitService struct {
	cache      database.Cache
	lock       *redsync.Redsync
	MaxRequest int
	Duration   time.Duration
}

func NewRateLimitService(maxRequest int, duration time.Duration, cache database.Cache, lock *redsync.Redsync) RateLimitService {
	return &rateLimitService{
		cache:      cache,
		MaxRequest: maxRequest,
		Duration:   duration,
		lock:       lock,
	}
}

func (s *rateLimitService) Allow(ctx context.Context, ipOrToken string) bool {

	if ipOrToken == "" {
		log.Fatalf("Failed to set key: %v", fmt.Errorf(ipOrToken+" is empty"))
	}

	mutex := s.lock.NewMutex("rateLimit")
	err := mutex.Lock()
	if err != nil {
		log.Fatalf("Failed to lock mutex: %v", err)
	}
	defer mutex.Unlock()

	newValue, err := s.cache.Incr(ctx, ipOrToken)
	if err != nil {
		log.Fatalf("Failed to increment key: %v", err)
	}

	if newValue == 1 {
		err := s.cache.Set(ctx, ipOrToken, 1, s.Duration)
		if err != nil {
			log.Fatalf("Failed to set key: %v", err)
		}
	} else {
		if newValue > int64(s.MaxRequest) {
			return false
		}
	}
	//exists, err := s.cache.Exists(ctx, ipOrToken)
	//if err != nil {
	//	log.Fatalf("Failed to check if key exists: %v", err)
	//}
	//if exists {
	//	newValue, err := s.cache.Incr(ctx, ipOrToken)
	//	if err != nil {
	//		log.Fatalf("Failed to increment key: %v", err)
	//	}
	//
	//	if newValue > int64(s.MaxRequest) {
	//		return false
	//	}
	//
	//} else {
	//	// first for this key, in the time window
	//	err = s.cache.Set(ctx, ipOrToken, 1, s.Duration)
	//	if err != nil {
	//		log.Fatalf("Failed to set key: %v", err)
	//	}
	//}

	return true
}
