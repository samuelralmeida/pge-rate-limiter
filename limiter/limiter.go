package limiter

import (
	"context"
	"fmt"
	"time"
)

type Limiter struct {
	storage      limiterStorage
	config       *limiterConfig
	tokenStorage tokenFetcher
}

type limiterStorage interface {
	Increment(ctx context.Context, key string, ttl time.Duration) (int, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Exists(ctx context.Context, key string) (int, error)
}

type tokenFetcher interface {
	GetLimitByToken(token string) int
}

func NewLimiter(storage limiterStorage, tokenStorage tokenFetcher) *Limiter {
	return &Limiter{
		storage:      storage,
		tokenStorage: tokenStorage,
		config:       config(),
	}
}

func (l *Limiter) isAllow(ctx context.Context, key string, maximum int) (bool, error) {
	total, err := l.storage.Increment(ctx, key, time.Second*60)
	if err != nil {
		return false, err
	}
	return total <= maximum, nil
}

func (l *Limiter) IsTokenAllow(ctx context.Context, token string) (bool, error) {
	total, err := l.storage.Exists(ctx, token)
	if err != nil {
		return false, fmt.Errorf("error to get ip status: %w", err)
	}

	if total > 0 {
		return false, nil
	}

	limit := l.tokenStorage.GetLimitByToken(token)
	if limit == 0 {
		return false, nil
	}

	allowed, err := l.isAllow(ctx, token, limit)
	if err != nil {
		return false, fmt.Errorf("error to check if token is allowed: %w", err)
	}

	if !allowed {
		l.storage.Set(ctx, token, "block", l.config.blockTokenDuration)
	}

	return allowed, nil
}

func (l *Limiter) IsIPAllow(ctx context.Context, ip string) (bool, error) {
	total, err := l.storage.Exists(ctx, ip)
	if err != nil {
		return false, fmt.Errorf("error to get ip status: %w", err)
	}

	if total > 0 {
		return false, nil
	}

	allowed, err := l.isAllow(ctx, ip, l.config.maxIPLimit)
	if err != nil {
		return false, fmt.Errorf("error to check if ip is allowed: %w", err)
	}

	if !allowed {
		l.storage.Set(ctx, ip, "block", l.config.blockIPDuration)
	}

	return allowed, nil
}
