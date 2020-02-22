package kubemqcluster

type RestConfig struct {
	// +optional
	Disabled bool `json:"disabled"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=(ClusterIP|NodePort|LoadBalancer)
	Expose string `json:"expose"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	BufferSize int32 `json:"bufferSize"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	BodyLimit int32 `json:"bodyLimit"`

	// +optional
	// +kubebuilder:validation:Minimum=30000
	// +kubebuilder:validation:Maximum=32767
	NodePort int32 `json:"nodePort"`
}

func (c *RestConfig) getDefaults() *RestConfig {
	if c.Port == 0 {
		c.Port = 9090
	}
	if c.Expose == "" {
		c.Expose = "ClusterIP"
	}
	return c
}
