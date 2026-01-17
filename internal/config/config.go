package config

import (
	"log"
	"os"
)

type Config struct {
	DBURL string
}

func Load() *Config {
	cfg := &Config{
		DBURL: os.Getenv("DB_URL"),
	}

	if cfg.DBURL == "" {
		log.Fatal("DB_URL is required")
	}

	return cfg
}
