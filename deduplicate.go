package main

import (
	"context"
	"golang.org/x/text/language"
	"sync"
)

type deduplicate struct {
	translator Translator
	dataMap    map[string]bool
	condition  *sync.Cond
}

func newDeduplicate(translator Translator) (Translator, error) {
	return &deduplicate{
		translator: translator,
		dataMap:    map[string]bool{},
		condition:  sync.NewCond(&sync.Mutex{}),
	}, nil
}

func (d *deduplicate) Translate(ctx context.Context, from, to language.Tag, data string) (string, error) {
	key := generateKey(from, to, data)

	d.condition.L.Lock()
	for d.isKeyExists(key) {
		d.condition.Wait()
	}
	d.dataMap[key] = true
	d.condition.L.Unlock()

	result, err := d.translator.Translate(ctx, from, to, data)
	if err != nil {
		return "", err
	}

	d.condition.L.Lock()
	delete(d.dataMap, key)
	d.condition.Broadcast()
	d.condition.L.Unlock()

	return result, nil
}

func (d *deduplicate) isKeyExists(key string) bool {
	_, ok := d.dataMap[key]

	return ok
}
