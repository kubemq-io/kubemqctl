package kubemqcluster

type HealthConfig struct {
	Enabled bool `json:"enabled"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	InitialDelaySeconds int32 `json:"initialDelaySeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	PeriodSeconds int32 `json:"periodSeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	TimeoutSeconds int32 `json:"timeoutSeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	SuccessThreshold int32 `json:"successThreshold"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	FailureThreshold int32 `json:"failureThreshold"`
}

func (c *HealthConfig) getDefaults() *HealthConfig {
	if c.InitialDelaySeconds == 0 {
		c.InitialDelaySeconds = 5
	}

	if c.PeriodSeconds == 0 {
		c.PeriodSeconds = 10
	}

	if c.TimeoutSeconds == 0 {
		c.TimeoutSeconds = 5
	}

	if c.SuccessThreshold == 0 {
		c.SuccessThreshold = 1
	}

	if c.FailureThreshold == 0 {
		c.FailureThreshold = 6
	}

	return c
}
