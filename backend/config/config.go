package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServiceModeType string

const (
	AuthoritativeDB = "AUTHORRITATIVE_DB"
	UpstreamCache   = "UPSTREAM_CACHE"
)

type Config struct {
	ServiceMode ServiceModeType
}

func defaultConfig() Config {
	return Config{
		ServiceMode: UpstreamCache,
	}
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Using default config no env file")
	}

	cfg := defaultConfig()

	if serviceMode, exists := os.LookupEnv("SERVICE_MODE"); exists {
		if serviceMode != UpstreamCache && serviceMode != AuthoritativeDB {
			return nil, errors.New("wrong SERVICE_MODE")
		}
		cfg.ServiceMode = ServiceModeType(serviceMode)
	}

	return &cfg, nil
}
