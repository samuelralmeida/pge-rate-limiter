package limiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/samuelralmeida/pge-rate-limiter/limiter"
)

type limiterMock struct {
	incrementErr    error
	incrementResult int
	setErr          error
	existsErr       error
	existsResp      int
}

func (lm *limiterMock) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	return lm.incrementResult, lm.incrementErr

}

func (lm *limiterMock) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return lm.setErr

}

func (lm *limiterMock) Exists(ctx context.Context, key string) (int, error) {
	return lm.existsResp, lm.existsErr
}

type tokenMock struct {
	limit int
}

func (tm *tokenMock) GetLimitByToken(token string) int {
	return tm.limit
}

func newConfig() *limiter.Config {
	return &limiter.Config{
		MaxIPLimit:         10,
		BlockIPDuration:    time.Duration(time.Second * 60),
		BlockTokenDuration: time.Duration(time.Second * 60),
	}
}

func TestLimiter_IsIPAllow_allow(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.IPMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 5,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      0,
	}

	limiter := limiter.NewLimiter(mock, nil, config)

	got, err := limiter.IsAllow(ctx, "123.456.789", "")

	var (
		wantErr  bool = false
		wantResp bool = true
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_IsIPAllow_notAllow(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.IPMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 11,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      0,
	}

	limiter := limiter.NewLimiter(mock, nil, config)

	got, err := limiter.IsAllow(ctx, "123.456.789", "")

	var (
		wantErr  bool = false
		wantResp bool = false
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_IsIPAllow_block(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.AnyMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 5,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      1,
	}

	tokenMock := &tokenMock{
		limit: 0,
	}

	limiter := limiter.NewLimiter(mock, tokenMock, config)

	got, err := limiter.IsAllow(ctx, "123.456.789", "")

	var (
		wantErr  bool = false
		wantResp bool = false
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_IsTokenAllow_allow(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.TokenMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 5,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      0,
	}

	tokenMock := &tokenMock{
		limit: 8,
	}

	limiter := limiter.NewLimiter(mock, tokenMock, config)

	got, err := limiter.IsAllow(ctx, "", "token-api")

	var (
		wantErr  bool = false
		wantResp bool = true
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_IsTokenAllow_notAllow(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.TokenMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 9,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      0,
	}

	tokenMock := &tokenMock{
		limit: 8,
	}

	limiter := limiter.NewLimiter(mock, tokenMock, config)

	got, err := limiter.IsAllow(ctx, "", "token-api")

	var (
		wantErr  bool = false
		wantResp bool = false
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_IsTokenAllow_block(t *testing.T) {
	ctx := context.Background()

	config := newConfig()
	config.Mode = limiter.AnyMode

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 5,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      1,
	}

	tokenMock := &tokenMock{
		limit: 8,
	}

	limiter := limiter.NewLimiter(mock, tokenMock, config)

	got, err := limiter.IsAllow(ctx, "", "token-api")

	var (
		wantErr  bool = false
		wantResp bool = false
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}

func TestLimiter_NoMode(t *testing.T) {
	ctx := context.Background()

	config := newConfig()

	mock := &limiterMock{
		incrementErr:    nil,
		incrementResult: 5,
		setErr:          nil,
		existsErr:       nil,
		existsResp:      1,
	}

	tokenMock := &tokenMock{
		limit: 8,
	}

	limiter := limiter.NewLimiter(mock, tokenMock, config)

	got, err := limiter.IsAllow(ctx, "123.456.789", "token-api")

	var (
		wantErr  bool = true
		wantResp bool = false
	)

	if (err != nil) != wantErr {
		t.Errorf("Limiter.IsIPAllow() error = %v, wantErr %v", err, wantErr)
		return
	}
	if got != wantResp {
		t.Errorf("Limiter.IsIPAllow() = %v, want %v", got, wantResp)
	}
}
