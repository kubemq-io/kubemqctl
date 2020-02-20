package cluster

type VolumeConfig struct {
	// +optional
	// +kubebuilder:validation:Pattern=^([1-9]?[0-9]?[0-9]?[0-9]?[0-9]?[0-9]?)Gi$
	Size string `json:"size"`
}
