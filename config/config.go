package config

import "os"

type Config struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
}

func New() *Config {
	return &Config{
		DbName:     os.Getenv("POSTGRES_DB"),
		DbUser:     os.Getenv("POSTGRES_USER"),
		DbPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbHost:     os.Getenv("DB_HOST"),
		DbPort:     os.Getenv("DB_PORT"),
	}
}
