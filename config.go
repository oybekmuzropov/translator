package main

import "os"

type config struct {

}

func loadConfig() *config {
	cfg := config{}

	return &cfg
}

func getOrReturnDefaultValue(key string, defaultValue interface{}) interface{} {
	_, ok := os.LookupEnv(key)
	if ok {
		return os.Getenv(key)
	}
	return defaultValue
}