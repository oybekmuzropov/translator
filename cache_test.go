package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
	"time"
)

func TestCache_TranslateHasCache(t *testing.T) {
	mockTranslator := mockTranslator{}

	cacheTranslator := cache{&mockTranslator, redisClient, time.Duration(1) * time.Minute}

	ctx := context.Background()
	from := language.Uzbek
	to := language.English
	data := "olma"
	expectedResult := "apple"

	err := cacheTranslator.set(ctx, from, to, data, expectedResult)
	assert.NoError(t, err)

	res, err := cacheTranslator.Translate(ctx, from, to, data)
	assert.NoError(t, err)
	assert.Equal(t, res, expectedResult)

	key := generateKey(from, to, data)
	val, err := redisClient.Get(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, val)

	reverseKey := generateKey(to, from, res)
	val, err = redisClient.Get(ctx, reverseKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, data, val)
}

func TestCache_TranslateNoCache(t *testing.T) {
	mockTranslator := mockTranslator{}

	cacheTranslator := cache{&mockTranslator, redisClient, time.Duration(1) * time.Minute}

	ctx := context.Background()
	from := language.Uzbek
	to := language.English
	data := "olma"
	expectedResult := "apple"

	mockTranslator.On("Translate", ctx, from, to, data).Return(expectedResult, nil)

	res, err := cacheTranslator.Translate(ctx, from, to, data)
	assert.NoError(t, err)
	assert.Equal(t, res, expectedResult)

	key := generateKey(from, to, data)
	val, err := redisClient.Get(ctx, key).Result()
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, val)

	reverseKey := generateKey(to, from, res)
	val, err = redisClient.Get(ctx, reverseKey).Result()
	assert.NoError(t, err)
	assert.Equal(t, data, val)
}