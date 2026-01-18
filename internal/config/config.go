package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBURL           string
	SecretKey       string
	Issuer          string
	Audience        string
	ExpirationHours int
}

func Load() *Config {
	expirationHours, err := strconv.Atoi(os.Getenv("EXPIRATION_HOURS"))
	if err != nil {
		log.Fatal("EXPIRATION_HOURS is required")
	}

	cfg := &Config{
		DBURL:           os.Getenv("DB_URL"),
		SecretKey:       os.Getenv("SECRET_KEY"),
		Issuer:          os.Getenv("ISSUER"),
		Audience:        os.Getenv("AUDIENCE"),
		ExpirationHours: expirationHours,
	}

	if cfg.DBURL == "" {
		log.Fatal("DB_URL is required")
	}

	if cfg.SecretKey == "" {
		log.Fatal("SECRET_KEY is required")
	}

	if cfg.Issuer == "" {
		log.Fatal("ISSUER is required")
	}

	if cfg.Audience == "" {
		log.Fatal("AUDIENCE is required")
	}

	return cfg
}
