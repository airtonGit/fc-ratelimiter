package ratelimiter

import (
	"context"
	"time"

	"ratelimiter-v2/internal/infrastructure/database"
)

type AllowRateLimiterUsecase interface {
	Execute(ctx context.Context, dto AllowRateLimitInputDTO) AllowRateLimitOutputDTO
}

type AllowRateLimitInputDTO struct {
	IpLimit    int
	IpDuration time.Duration

	TokenLimit    int
	TokenDuration time.Duration

	IP    string
	Token string
}

type AllowRateLimitOutputDTO struct {
	Allow bool
}

func NewRateLimiterUsecase(cache database.Cache) AllowRateLimiterUsecase {
	return &rateLimiterUsecaseImpl{
		cache: cache,
	}
}

type rateLimiterUsecaseImpl struct {
	cache database.Cache
}

func (r *rateLimiterUsecaseImpl) Execute(ctx context.Context, dto AllowRateLimitInputDTO) AllowRateLimitOutputDTO {
	ipRateLimiter := NewRateLimitService(dto.IpLimit, dto.IpDuration, r.cache)
	tokenRateLimiter := NewRateLimitService(dto.TokenLimit, dto.TokenDuration, r.cache)

	ipAllowed := ipRateLimiter.Allow(ctx, dto.IP)
	tokenAllowed := tokenRateLimiter.Allow(ctx, dto.Token)

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
