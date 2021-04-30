package main

import (
	"fmt"
	"os"
	"time"
)

// Service is a Translator user.
type Service struct {
	translator Translator
}

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)
	config := loadConfig()

	retryTranslator, err := newRetry(t, config)
	if err != nil {
		fmt.Printf("retryTranslator failed: %+v\n\n", err)
		os.Exit(1)
	}

	cacheTranslator, err := newCache(retryTranslator, config)
	if err != nil {
		fmt.Printf("cacheTranslator failed: %+v\n\n", err)
		os.Exit(1)
	}

	return &Service{
		translator: cacheTranslator,
	}
}
