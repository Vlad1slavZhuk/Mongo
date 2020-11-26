package config

import "os"

type Config struct {
	PortHTTP string
}

func NewConfigFromEnv() *Config {
	var ok bool
	config := new(Config)
	if config.PortHTTP, ok = os.LookupEnv("SERVER_PORT"); !ok {
		config.PortHTTP = "8000"
	}
	config.PortHTTP = ":" + config.PortHTTP

	return config
}
