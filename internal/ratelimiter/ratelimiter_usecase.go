package ratelimiter

import (
	"context"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
)

type AllowRateLimiterUsecase interface {
	Execute(ctx context.Context, dto AllowRateLimitInputDTO) AllowRateLimitOutputDTO
}

type AllowRateLimitInputDTO struct {
	IPRequestBySecondLimit int
	IpDuration             time.Duration

	TokenRequestsBySecondLimit map[string]int
	TokenDuration              time.Duration

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
	// validation
	if dto.IP == "" && dto.Token == "" {
		return AllowRateLimitOutputDTO{
			Allow: false,
		}
	}

	ipRateLimiter := NewRateLimitService(dto.IPRequestBySecondLimit, dto.IpDuration, r.cache)
	tokenRateLimiter := NewRateLimitService(dto.TokenRequestsBySecondLimit[dto.Token], dto.TokenDuration, r.cache)

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
