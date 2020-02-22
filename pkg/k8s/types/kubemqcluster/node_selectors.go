package kubemqcluster

type NodeSelectorConfig struct {
	Keys map[string]string `json:"keys"`
}

func (c *NodeSelectorConfig) DeepCopy() *NodeSelectorConfig {
	out := &NodeSelectorConfig{
		Keys: map[string]string{},
	}
	if out.Keys == nil {
		out.Keys = map[string]string{}
	}
	for key, value := range c.Keys {
		out.Keys[key] = value
	}
	return out
}
