package main

import "os"

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getDBPath() string {
	env := getEnv("RENDER", "")
	if env != "" {
		return "/tmp/site.db"
	}
	return getEnv("DB_PATH", "./site.db")
}
