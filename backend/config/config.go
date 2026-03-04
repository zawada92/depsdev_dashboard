package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type ServiceModeType string

const (
	AuthoritativeDB = "AUTHORITATIVE_DB"
	UpstreamCache   = "UPSTREAM_CACHE"
)

type Config struct {
	ServiceMode          ServiceModeType
	HttpClientTimeoutSec int
	DbPath               string
	CorsAddress          string
}

func defaultConfig() Config {
	return Config{
		ServiceMode:          UpstreamCache,
		HttpClientTimeoutSec: 10,
		DbPath:               "/data/app.db",
		CorsAddress:          "*",
	}
}

// TODO use viper or something
func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("Using default config no env file")
	}

	cfg := defaultConfig()

	if serviceMode, exists := os.LookupEnv("SERVICE_MODE"); exists {
		if serviceMode != UpstreamCache && serviceMode != AuthoritativeDB {
			return nil, errors.New("wrong SERVICE_MODE")
		}
		cfg.ServiceMode = ServiceModeType(serviceMode)
	}
	if httpClientTimeout, exists := os.LookupEnv("HTTP_CLIENT_TIMEOUT_SEC"); exists {
		timeout, err := strconv.Atoi(httpClientTimeout)
		if err != nil || timeout <= 0 {
			return nil, errors.New("wrong HTTP_CLIENT_TIMEOUT_SEC")
		}
		cfg.HttpClientTimeoutSec = timeout
	}
	if dbPath, exists := os.LookupEnv("DB_PATH"); exists {
		cfg.DbPath = dbPath
	}
	if corsAddress, exists := os.LookupEnv("CORS_ADDRESS"); exists {
		cfg.CorsAddress = corsAddress
	}

	return &cfg, nil
}
