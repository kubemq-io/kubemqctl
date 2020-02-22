package kubemqcluster

import (
	"os"
)

type LicenseConfig struct {

	// +optional
	Data string `json:"data"`

	// +optional
	Token string `json:"token"`
}

func (c *LicenseConfig) getFromEnv(currentValue string, envKey string, def string) string {
	if currentValue != "" {
		return currentValue
	}
	fromEnv := os.Getenv(envKey)
	if fromEnv != "" {
		return fromEnv
	}
	return def
}
