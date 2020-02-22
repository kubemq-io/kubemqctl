package kubemqcluster

type ApiConfig struct {
	// +optional
	Disabled bool `json:"disabled"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// +optional
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose"`

	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	NodePort int32 `json:"nodePort"`
}

func (c *ApiConfig) getDefaults() *ApiConfig {
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Expose == "" {
		c.Expose = "ClusterIP"
	}
	return c
}
