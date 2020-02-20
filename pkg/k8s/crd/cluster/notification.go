package cluster

type NotificationConfig struct {
	Enabled bool `json:"enabled"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Prefix string `json:"prefix"`

	// +optional
	Log bool `json:"log"`
}
