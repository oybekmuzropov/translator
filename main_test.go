package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"os"
	"testing"
)

var (
	redisCache *redis.Client
	cfg *config
)

func TestMain(m *testing.M)  {
	cfg = loadConfig()
	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.redisHost, cfg.redisPort),
			Password: cfg.redisPassword,
		})
	})

	code := m.Run()

	os.Exit(code)
}