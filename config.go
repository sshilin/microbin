package main

import (
	"os"
	"strings"
)

type Config struct {
	Listen        string
	TLS           Tls
	LogFormatJson bool
}

type Tls struct {
	Enabled  bool
	KeyFile  string
	CertFile string
}

func LoadConfig() Config {
	return Config{
		Listen:        lookupEnv("LISTEN", "0.0.0.0:8080"),
		LogFormatJson: lookupEnvBool("LOG_FORMAT_JSON", "false"),
		TLS: Tls{
			Enabled:  lookupEnvBool("TLS_ENABLED", "false"),
			KeyFile:  lookupEnv("TLS_KEY_FILE", ""),
			CertFile: lookupEnv("TLS_CERT_FILE", ""),
		},
	}
}

func lookupEnv(key, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}

func lookupEnvBool(key, defaultValue string) bool {
	return strings.EqualFold(lookupEnv(key, defaultValue), "true")
}
