package api

import "register/storage"

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LoggerLevel string `toml:"logger_level"`
	Storage     *storage.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8000",
		LoggerLevel: "debug",
		Storage:     storage.NewConfig(),
	}
}
