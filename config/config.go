package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	port string
}

func (conf *Config) Port() string {
	return conf.port
}

func Load() *Config {
    // Load the default .env file into the environment if it exists
	godotenv.Load()

	conf := &Config{
		port: os.Getenv("PORT"),
	}

	return conf
}
