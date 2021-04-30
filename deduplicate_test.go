package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"sync"
	"testing"
	"time"
)

type mockDeduplicateTranslator struct {
	mock.Mock
}

func (t *mockDeduplicateTranslator) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	args := t.Called(ctx, from, to, data)
	time.Sleep(time.Duration(1) * time.Second)
	return args.String(0), args.Error(1)
}

func TestDeduplicate_Translate(t *testing.T) {
	mockTranslator := mockDeduplicateTranslator{}
	translator, err := newDeduplicate(&mockTranslator)
	assert.NoError(t, err)

	ctx := context.Background()
	from := language.Uzbek
	to := language.Uzbek
	data := "olma"
	expectedResult := "apple"

	mockTranslator.On("Translate", ctx, from, to, data).Return(expectedResult, nil)

	var wg sync.WaitGroup

	for i := 1; i < 10; i ++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := translator.Translate(ctx, from, to, data)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}