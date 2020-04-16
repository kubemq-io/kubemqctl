package kubemqcluster

type QueueConfig struct {

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxReceiveMessagesRequest *int32 `json:"maxReceiveMessagesRequest"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxWaitTimeoutSeconds *int32 `json:"maxWaitTimeoutSeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxExpirationSeconds *int32 `json:"maxExpirationSeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxDelaySeconds *int32 `json:"maxDelaySeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxReQueues *int32 `json:"maxReQueues"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	MaxVisibilitySeconds *int32 `json:"maxVisibilitySeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultVisibilitySeconds *int32 `json:"defaultVisibilitySeconds"`

	// +optional
	// +kubebuilder:validation:Minimum=1
	DefaultWaitTimeoutSeconds *int32 `json:"defaultWaitTimeoutSeconds"`
}

func (c *QueueConfig) DeepCopy() *QueueConfig {
	out := &QueueConfig{}

	if c.MaxReceiveMessagesRequest != nil {
		out.MaxReceiveMessagesRequest = new(int32)
		*out.MaxReceiveMessagesRequest = *c.MaxReceiveMessagesRequest
	}

	if c.MaxWaitTimeoutSeconds != nil {
		out.MaxWaitTimeoutSeconds = new(int32)
		*out.MaxWaitTimeoutSeconds = *c.MaxWaitTimeoutSeconds

	}

	if c.MaxExpirationSeconds != nil {
		out.MaxExpirationSeconds = new(int32)
		*out.MaxExpirationSeconds = *c.MaxExpirationSeconds
	}

	if c.MaxDelaySeconds != nil {
		out.MaxDelaySeconds = new(int32)
		*out.MaxDelaySeconds = *c.MaxDelaySeconds
	}

	if c.MaxReQueues != nil {
		out.MaxReQueues = new(int32)
		*out.MaxReQueues = *c.MaxReQueues
	}

	if c.MaxVisibilitySeconds != nil {
		out.MaxVisibilitySeconds = new(int32)
		*out.MaxVisibilitySeconds = *c.MaxVisibilitySeconds
	}

	if c.DefaultVisibilitySeconds != nil {
		out.DefaultVisibilitySeconds = new(int32)
		*out.DefaultVisibilitySeconds = *c.DefaultVisibilitySeconds
	}

	if c.DefaultWaitTimeoutSeconds != nil {
		out.DefaultWaitTimeoutSeconds = new(int32)
		*out.DefaultWaitTimeoutSeconds = *c.DefaultWaitTimeoutSeconds
	}

	return out
}
