package main

import (
	"github.com/spf13/cast"
	"os"
)

type config struct {
	retryTimes int

	redisHost          string
	redisPort          int
	redisPassword      string
	redisCacheDuration int // hour
}

func loadConfig() *config {
	cfg := config{}

	cfg.retryTimes = cast.ToInt(getOrReturnDefaultValue("RETRY_TIMES", 3))

	cfg.redisHost = cast.ToString(getOrReturnDefaultValue("REDIS_HOST", "localhost"))
	cfg.redisPort = cast.ToInt(getOrReturnDefaultValue("REDIS_PORT", 6379))
	cfg.redisPassword = cast.ToString(getOrReturnDefaultValue("REDIS_PASSWORD", ""))
	cfg.redisCacheDuration = cast.ToInt(getOrReturnDefaultValue("REDIS_CACHE_DURATION", 1))

	return &cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, ok := os.LookupEnv(key)
	if ok {
		return os.Getenv(key)
	}
	return defaultValue
}
