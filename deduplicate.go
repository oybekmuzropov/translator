package main

import (
	"context"
	"golang.org/x/text/language"
)

type deduplicate struct {
	translator Translator
}

func newDeduplicate(translator Translator) (Translator, error) {
	return &deduplicate{
		translator: translator,
	}, nil
}

func (d *deduplicate) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	// TODO should implement
	return "", nil
}