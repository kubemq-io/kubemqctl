package config

type KubernetesConfig struct {
	AutoIntegrated bool
	Namespace      string
	StatefulSet    string
	Pod            string
}
