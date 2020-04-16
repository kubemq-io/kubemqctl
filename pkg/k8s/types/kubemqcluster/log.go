package kubemqcluster

type LogConfig struct {
	// +optional
	// +kubebuilder:validation:Minimum=0
	Level *int32 `json:"level"`

	// +optional
	// +kubebuilder:validation:MinLength=1
	File string `json:"file"`
}

func (c *LogConfig) DeepCopy() *LogConfig {
	out := &LogConfig{}
	out.File = c.File
	if c.Level != nil {
		out.Level = new(int32)
		*out.Level = *c.Level
	}
	return out
}
