package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"testing"
)

type mockTranslator struct {
	mock.Mock
}

var (
	wrongTranslation = errors.New("translation is incorrect")
)

func (t *mockTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	args := t.Called(ctx, from, to, data)

	return args.String(0), args.Error(1)
}

func TestRetry_TranslateWithNoError(t *testing.T) {
	mockTranslator := mockTranslator{}
	retryTranslator := &retry{
		translator: &mockTranslator,
		config:     cfg,
	}

	ctx := context.Background()
	from := language.Uzbek
	to := language.English
	data := "olma"
	expectedResult := "apple"

	mockTranslator.On("Translate", ctx, from, to, data).Return(expectedResult, nil)

	result, err := retryTranslator.Translate(ctx, from, to, data)
	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestRetry_TranslateWithError(t *testing.T) {
	mockTranslator := mockTranslator{}
	cfg.retryTimes = 1
	retryTranslator := &retry{
		translator: &mockTranslator,
		config:     cfg,
	}

	ctx := context.Background()
	from := language.Uzbek
	to := language.English
	data := "olma"

	mockTranslator.On("Translate", ctx, from, to, data).Return("", wrongTranslation)

	result, err := retryTranslator.Translate(ctx, from, to, data)
	assert.ErrorIs(t, err, wrongTranslation)
	assert.Empty(t, result)
}