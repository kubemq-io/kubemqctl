package cluster

type ResourceConfig struct {

	// +kubebuilder:validation:MinLength=1
	LimitsCpu string `json:"limitsCpu"`

	// +kubebuilder:validation:MinLength=1
	LimitsMemory string `json:"limitsMemory"`

	// +kubebuilder:validation:MinLength=1
	RequestsCpu string `json:"requestsCpu"`

	// +kubebuilder:validation:MinLength=1
	RequestsMemory string `json:"requestsMemory"`
}
