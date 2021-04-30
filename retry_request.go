package main

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"time"
)

type retry struct {
	translator Translator
	config     *config
}

func newRetry(t Translator, config *config) (Translator, error) {
	if config.retryTimes <= 0 {
		return nil, errors.New("retry times must be bigger than 0")
	}
	return &retry{
		translator: t,
		config:     config,
	}, nil
}

func (b *retry) Translate(ctx context.Context, from, to language.Tag, data string) (result string, err error) {
	for i := 1; i <= b.config.retryTimes; i++ {
		result, err = b.translator.Translate(ctx, from, to, data)
		if err == nil {
			return result, nil
		}

		time.Sleep(time.Duration(i*i) * time.Second)
	}
	return "", errors.Wrap(err, "failed while translating")
}
