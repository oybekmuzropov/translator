package main

import (
	"os"
	"github.com/spf13/cast"
)

type config struct {
	retryTimes int
}

func loadConfig() *config {
	cfg := config{}

	cfg.retryTimes = cast.ToInt(getOrReturnDefaultValue("RETRY_TIMES", -3))

	return &cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, ok := os.LookupEnv(key)
	if ok {
		return os.Getenv(key)
	}
	return defaultValue
}