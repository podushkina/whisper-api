package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Server struct {
		Port string `env:"SERVER_PORT" env-default:"8080"`
	}
	Redis struct {
		Addr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
		Password string `env:"REDIS_PASSWORD" env-default:""`
		DB       int    `env:"REDIS_DB" env-default:"0"`
	}
	JWT struct {
		Secret string `env:"JWT_SECRET" env-required:"true"`
	}
	Whisper struct {
		Path string `env:"WHISPER_PATH" env-required:"true"`
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Printf("Error reading environment variables: %v", err)
		return nil, err
	}

	log.Printf("Configuration loaded successfully: %+v", cfg)

	return cfg, nil
}
