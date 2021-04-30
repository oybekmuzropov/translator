package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"sync"
	"time"
)

var (
	once        sync.Once
	redisClient *redis.Client
)

type cache struct {
	translator Translator
	cache      *redis.Client
	duration   time.Duration
}

func generateKey(from, to language.Tag, data string) string {
	return fmt.Sprintf("%s-%s-%s", from, to, data)
}

func newCache(translator Translator, config *config) (Translator, error) {
	if config.redisCacheDuration <= 0 {
		return nil, errors.New("redis cache duration must be bigger than 0")
	}

	once.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.redisHost, config.redisPort),
			Password: config.redisPassword,
		})
	})

	return &cache{
		translator: translator,
		cache:      redisClient,
		duration:   time.Duration(config.redisCacheDuration) * time.Hour,
	}, nil
}

func (c *cache) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := generateKey(from, to, data)

	val, err := c.cache.Get(ctx, key).Result()
	if err == redis.Nil {
		result, err := c.translator.Translate(ctx, from, to, data)
		if err != nil {
			return "", err
		}

		err = c.cache.Set(ctx, key, result, c.duration).Err()
		if err != nil {
			return "", err
		}

		reverseKey := generateKey(to, from, result)
		err = c.cache.Set(ctx, reverseKey, data, c.duration).Err()
		if err != nil {
			return "", err
		}

		return result, nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}
