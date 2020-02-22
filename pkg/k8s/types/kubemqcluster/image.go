package kubemqcluster

import (
	"fmt"
	"os"
)

type ImageConfig struct {

	// +optional
	Registry string `json:"registry"`

	// +optional
	Repository string `json:"repository"`

	// +optional
	Tag string `json:"tag"`

	// +optional
	// +kubebuilder:validation:Pattern=(IfNotPresent|Always|Never)
	PullPolicy string `json:"pullPolicy"`
}

func (c *ImageConfig) getFromEnv(currentValue string, envKey string, def string) string {
	if currentValue != "" {
		return currentValue
	}
	fromEnv := os.Getenv(envKey)
	if fromEnv != "" {
		return fromEnv
	}
	return def
}
func (c *ImageConfig) GetImage() string {
	return fmt.Sprintf("%s/%s:%s", c.Registry, c.Repository, c.Tag)
}
