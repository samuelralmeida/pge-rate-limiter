package limiter

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type mode string

const (
	IPMode    mode = "ip-mode"
	TokenMode mode = "token-mode"
	AnyMode   mode = "any-mode"
)

type Config struct {
	MaxIPLimit         int
	BlockIPDuration    time.Duration
	BlockTokenDuration time.Duration
	Mode               mode
}

func NewConfig() *Config {
	maxIPLimit, err := strconv.Atoi(os.Getenv("MAX_IP_LIMIT"))
	if err != nil {
		maxIPLimit = 10
	}

	blockIpSeconds, err := strconv.Atoi(os.Getenv("BLOCK_IP_SECONDS"))
	if err != nil {
		blockIpSeconds = 60 * 5
	}

	blockTokenSeconds, err := strconv.Atoi(os.Getenv("BLOCK_TOKEN_SECONDS"))
	if err != nil {
		blockTokenSeconds = 60 * 5
	}

	modeLimiter := AnyMode
	mode := strings.ToUpper(os.Getenv("LIMITER_MODE"))
	if mode == "IP" {
		modeLimiter = IPMode
	} else if mode == "TOKEN" {
		modeLimiter = TokenMode
	}

	return &Config{
		MaxIPLimit:         maxIPLimit,
		BlockTokenDuration: time.Duration(time.Second * time.Duration(blockTokenSeconds)),
		BlockIPDuration:    time.Duration(time.Second * time.Duration(blockIpSeconds)),
		Mode:               modeLimiter,
	}
}
