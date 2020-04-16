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
