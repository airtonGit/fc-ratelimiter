package ratelimiter

import (
	"context"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/go-redsync/redsync/v4"
)

type AllowRateLimiterUsecase interface {
	Execute(ctx context.Context, dto AllowRateLimitInputDTO) AllowRateLimitOutputDTO
}

type AllowRateLimitInputDTO struct {
	IpLimit    int
	IpDuration time.Duration

	TokenLimit    map[string]int
	TokenDuration time.Duration

	IP    string
	Token string
}

type AllowRateLimitOutputDTO struct {
	Allow bool
}

func NewRateLimiterUsecase(cache database.Cache, lock *redsync.Redsync) AllowRateLimiterUsecase {
	return &rateLimiterUsecaseImpl{
		cache: cache,
		lock:  lock,
	}
}

type rateLimiterUsecaseImpl struct {
	cache database.Cache
	lock  *redsync.Redsync
}

func (r *rateLimiterUsecaseImpl) Execute(ctx context.Context, dto AllowRateLimitInputDTO) AllowRateLimitOutputDTO {
	// validation
	if dto.IP == "" && dto.Token == "" {
		return AllowRateLimitOutputDTO{
			Allow: false,
		}
	}

	ipRateLimiter := NewRateLimitService(dto.IpLimit, dto.IpDuration, r.cache, r.lock)
	tokenRateLimiter := NewRateLimitService(dto.TokenLimit[dto.Token], dto.TokenDuration, r.cache, r.lock)

	ipAllowed := ipRateLimiter.Allow(ctx, dto.IP)
	var tokenAllowed bool
	if dto.Token != "" {
		tokenAllowed = tokenRateLimiter.Allow(ctx, dto.Token)
	}

	if dto.Token != "" {
		return AllowRateLimitOutputDTO{
			Allow: tokenAllowed,
		}
	} else {
		return AllowRateLimitOutputDTO{
			Allow: ipAllowed,
		}
	}
}
