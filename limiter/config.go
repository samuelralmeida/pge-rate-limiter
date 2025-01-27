package limiter

import (
	"os"
	"strconv"
	"time"
)

type limiterConfig struct {
	maxIPLimit         int
	maxTokenLimit      int
	blockIPDuration    time.Duration
	blockTokenDuration time.Duration
}

func config() *limiterConfig {
	maxIPLimit, err := strconv.Atoi(os.Getenv("MAX_IP_LIMIT"))
	if err != nil {
		maxIPLimit = 10
	}

	maxTokenLimit, err := strconv.Atoi(os.Getenv("MAX_TOKEN_LIMIT"))
	if err != nil {
		maxTokenLimit = 10
	}

	blockIpSeconds, err := strconv.Atoi(os.Getenv("BLOCK_IP_SECONDS"))
	if err != nil {
		blockIpSeconds = 60 * 5
	}

	blockTokenSeconds, err := strconv.Atoi(os.Getenv("BLOCK_TOKEN_SECONDS"))
	if err != nil {
		blockTokenSeconds = 60 * 5
	}

	return &limiterConfig{
		maxIPLimit:         maxIPLimit,
		maxTokenLimit:      maxTokenLimit,
		blockTokenDuration: time.Duration(time.Second * time.Duration(blockTokenSeconds)),
		blockIPDuration:    time.Duration(time.Second * time.Duration(blockIpSeconds)),
	}
}
