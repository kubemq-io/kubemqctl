package deployment

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	"strings"
)

var defaultKubeMQConfigMapTemplate = `
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app.kubernetes.io/name: {{.Name}}
    app.kubernetes.io/instance: {{.Name}}
    app.kubernetes.io/managed-by: kubemqctl
    deployment.id: {{.Id}}
data:
  CLUSTER_ENABLE: "true"
{{ range $key, $value := .Variables}}
  {{$key}}: "{{$value}}"
{{end}}
`

type ConfigMapConfig struct {
	Id        string
	Name      string
	Namespace string
	Variables map[string]string
	configMap *apiv1.ConfigMap
}

func ImportConfigMap(spec []byte) (*ConfigMapConfig, error) {
	cm := &apiv1.ConfigMap{}
	err := yaml.Unmarshal(spec, cm)
	if err != nil {
		return nil, err
	}
	return &ConfigMapConfig{
		Id:        "",
		Name:      cm.Name,
		Namespace: cm.Namespace,
		Variables: nil,
		configMap: cm,
	}, nil
}

func NewConfigMap(id, name, namespace string) *ConfigMapConfig {
	return &ConfigMapConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		Variables: map[string]string{},
	}
}
func DefaultConfigMap(id, name, namespace string) map[string]*ConfigMapConfig {
	cm := make(map[string]*ConfigMapConfig)
	cm[name] = &ConfigMapConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		Variables: map[string]string{},
	}
	return cm
}

func (c *ConfigMapConfig) SetVariable(key, value string) *ConfigMapConfig {
	c.Variables[strings.ToUpper(key)] = value
	return c
}
func (c *ConfigMapConfig) Spec() ([]byte, error) {
	t := NewTemplate(defaultKubeMQConfigMapTemplate, c)
	return t.Get()
}
func (c *ConfigMapConfig) Set(value *apiv1.ConfigMap) *ConfigMapConfig {
	c.configMap = value
	return c
}
func (c *ConfigMapConfig) Get() (*apiv1.ConfigMap, error) {
	if c.configMap != nil {
		return c.configMap, nil
	}
	cm := &apiv1.ConfigMap{}
	data, err := c.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, cm)
	if err != nil {
		return nil, err
	}
	c.configMap = cm
	return cm, nil
}
