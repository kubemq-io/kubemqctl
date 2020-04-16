package kubemqcluster

type ImageConfig struct {

	// +optional
	Image string `json:"image"`

	// +optional
	// +kubebuilder:validation:Pattern=(IfNotPresent|Always|Never)
	PullPolicy string `json:"pullPolicy"`
}
