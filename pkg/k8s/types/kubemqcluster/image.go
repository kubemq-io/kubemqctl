package kubemqcluster

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
