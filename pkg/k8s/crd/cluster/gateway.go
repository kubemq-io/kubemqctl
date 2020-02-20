package cluster

type GatewayConfig struct {

	// +kubebuilder:validation:MinItems=1
	Remotes []string `json:"remotes"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Cert string `json:"cert"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Ca string `json:"ca"`
}

func (c *GatewayConfig) DeepCopy() *GatewayConfig {
	out := &GatewayConfig{
		Remotes: []string{},
		Cert:    c.Cert,
		Key:     c.Key,
		Ca:      c.Ca,
		Port:    c.Port,
	}
	for i := 0; i < len(c.Remotes); i++ {
		out.Remotes = append(out.Remotes, c.Remotes[i])
	}
	return out
}
