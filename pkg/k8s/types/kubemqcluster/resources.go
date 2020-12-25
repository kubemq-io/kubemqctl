package kubemqcluster

type ResourceConfig struct {

	// +optional
	LimitsCpu string `json:"limitsCpu,omitempty"`
	// +optional
	LimitsMemory string `json:"limitsMemory,omitempty"`
	// +optional
	LimitsEphemeralStorage string `json:"limitsEphemeralStorage,omitempty"`

	// +optional
	RequestsCpu string `json:"requestsCpu,omitempty"`
	// +optional
	RequestsMemory string `json:"requestsMemory,omitempty"`

	// +optional
	RequestsEphemeralStorage string `json:"requestsEphemeralStorage,omitempty"`
}
