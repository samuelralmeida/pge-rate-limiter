package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	MaxIPLimit         int
	MaxTokenLimit      int
	BlockIPDuration    time.Duration
	BlockTokenDuration time.Duration
}

func NewConfig() *Config {
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

	return &Config{
		MaxIPLimit:         maxIPLimit,
		MaxTokenLimit:      maxTokenLimit,
		BlockTokenDuration: time.Duration(time.Second * time.Duration(blockTokenSeconds)),
		BlockIPDuration:    time.Duration(time.Second * time.Duration(blockIpSeconds)),
	}
}
