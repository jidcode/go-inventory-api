package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Variables struct {
	DatabaseUrl string
	RedisUrl    string
	JWTSecret   string
}

func LoadEnv() *Variables {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load environmental variables: %s", err)
	}

	config := &Variables{
		DatabaseUrl: os.Getenv("DATABASE_URL"),
		RedisUrl:    os.Getenv("REDIS_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	return config
}
