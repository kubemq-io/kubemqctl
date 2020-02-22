package kubemqcluster

type StoreConfig struct {
	// +optional
	Clean bool `json:"clean"`

	// +optional
	// +kubebuilder:validation:MinLength=0
	Path string `json:"path"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannels *int32 `json:"maxChannels"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxSubscribers *int32 `json:"maxSubscribers"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxMessages *int32 `json:"maxMessages"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MaxChannelSize *int32 `json:"maxChannelSize"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	MessagesRetentionMinutes *int32 `json:"messagesRetentionMinutes"`

	// +optional
	// +kubebuilder:validation:Minimum=0
	PurgeInactiveMinutes *int32 `json:"purgeInactiveMinutes"`
}

func (c *StoreConfig) DeepCopy() *StoreConfig {
	out := &StoreConfig{}

	out.Clean = c.Clean
	out.Path = c.Path

	if c.MaxChannels != nil {
		out.MaxChannels = new(int32)
		*out.MaxChannels = *c.MaxChannels
	}

	if c.MaxSubscribers != nil {
		out.MaxSubscribers = new(int32)
		*out.MaxSubscribers = *c.MaxSubscribers

	}

	if c.MaxMessages != nil {
		out.MaxMessages = new(int32)
		*out.MaxMessages = *c.MaxMessages
	}

	if c.MaxChannelSize != nil {
		out.MaxChannelSize = new(int32)
		*out.MaxChannelSize = *c.MaxChannelSize
	}

	if c.MessagesRetentionMinutes != nil {
		out.MessagesRetentionMinutes = new(int32)
		*out.MessagesRetentionMinutes = *c.MessagesRetentionMinutes
	}

	if c.PurgeInactiveMinutes != nil {
		out.PurgeInactiveMinutes = new(int32)
		*out.PurgeInactiveMinutes = *c.PurgeInactiveMinutes
	}

	return out
}
