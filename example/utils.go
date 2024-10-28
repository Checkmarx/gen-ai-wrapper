package main

import (
	"fmt"
	"os"
)

func GetEnvKeyValue(key string) (string, error) {
	value, err := envLookup(key, true, "")
	if err != nil {
		return "", err
	}
	return value, nil
}

func envLookup(key string, mandatory bool, defaultValue string) (string, error) {
	val, found := os.LookupEnv(key)
	if !found {
		if !mandatory {
			return defaultValue, nil
		} else {
			return "", fmt.Errorf("environment variable '%s' not found", key)
		}
	}
	return val, nil
}
