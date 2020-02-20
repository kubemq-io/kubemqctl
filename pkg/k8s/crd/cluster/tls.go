package cluster

type TlsConfig struct {
	// +kubebuilder:validation:MinLength=1
	Cert string `json:"cert"`

	// +kubebuilder:validation:MinLength=1
	Key string `json:"key"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	Ca string `json:"ca"`
}
