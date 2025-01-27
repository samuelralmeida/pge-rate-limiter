package limiter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

type Limiter struct {
	storage      limiterStorage
	config       *Config
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

func NewLimiter(storage limiterStorage, tokenStorage tokenFetcher, config *Config) *Limiter {
	return &Limiter{
		storage:      storage,
		tokenStorage: tokenStorage,
		config:       config,
	}
}

func (l *Limiter) IsAllow(ctx context.Context, ip, token string) (bool, error) {
	if l.config.Mode == IPMode {
		return l.isIPAllow(ctx, ip)
	}

	if l.config.Mode == TokenMode {
		return l.isTokenAllow(ctx, token)
	}

	if l.config.Mode == AnyMode {
		if token != "" {
			return l.isTokenAllow(ctx, token)
		}
		return l.isIPAllow(ctx, ip)
	}

	return false, errors.New("rate limiter mode invalid")

}

func (l *Limiter) isAllow(ctx context.Context, key string, maximum int) (bool, error) {
	total, err := l.storage.Increment(ctx, key, time.Second*60)
	if err != nil {
		return false, err
	}
	return total <= maximum, nil
}

func (l *Limiter) isTokenAllow(ctx context.Context, token string) (bool, error) {
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
		err := l.storage.Set(ctx, token, "block", l.config.BlockTokenDuration)
		if err != nil {
			log.Println("error to block token: %w", err)
			return allowed, nil
		}
	}

	return allowed, nil
}

func (l *Limiter) isIPAllow(ctx context.Context, ip string) (bool, error) {
	total, err := l.storage.Exists(ctx, ip)
	if err != nil {
		return false, fmt.Errorf("error to get ip status: %w", err)
	}

	if total > 0 {
		return false, nil
	}

	allowed, err := l.isAllow(ctx, ip, l.config.MaxIPLimit)
	if err != nil {
		return false, fmt.Errorf("error to check if ip is allowed: %w", err)
	}

	if !allowed {
		err = l.storage.Set(ctx, ip, "block", l.config.BlockIPDuration)
		if err != nil {
			log.Println("error to block ip: %w", err)
			return allowed, nil
		}
	}

	return allowed, nil
}
