package config

import (
	"fmt"
	"os"
)

func Env(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return val
}
