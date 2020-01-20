package deployment

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	"strings"
)

var defaultKubeMQSecretTemplate = `
apiVersion: v1
kind: Secret
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.Name}}
    deployment.id: {{.Id}}
type: Opaque
stringData:
{{ range $key, $value := .Variables}}
  {{$key}}: "{{$value}}"
{{end}}
`

type SecretConfig struct {
	Id        string
	Name      string
	Namespace string
	Variables map[string]string
	secret    *apiv1.Secret
}

func ImportSecret(spec []byte) (*SecretConfig, error) {
	sec := &apiv1.Secret{}
	err := yaml.Unmarshal(spec, sec)
	if err != nil {
		return nil, err
	}
	return &SecretConfig{
		Id:        "",
		Name:      sec.Name,
		Namespace: sec.Namespace,
		Variables: nil,
		secret:    sec,
	}, nil
}
func NewSecretConfig(id, name, namespace string) *SecretConfig {
	return &SecretConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		Variables: map[string]string{},
	}
}
func DefaultSecretConfig(id, name, namespace string) map[string]*SecretConfig {
	secs := make(map[string]*SecretConfig)
	secs[name] = &SecretConfig{
		Id:        id,
		Name:      name,
		Namespace: namespace,
		Variables: map[string]string{},
	}
	return secs
}
func (s *SecretConfig) SetVariable(key, value string) *SecretConfig {
	s.Variables[strings.ToUpper(key)] = value
	return s
}
func (s *SecretConfig) Spec() ([]byte, error) {
	if s.secret == nil {
		t := NewTemplate(defaultKubeMQSecretTemplate, s)
		return t.Get()
	}
	return yaml.Marshal(s.secret)
}

func (s *SecretConfig) Set(value *apiv1.Secret) *SecretConfig {
	s.secret = value
	return s
}
func (s *SecretConfig) Get() (*apiv1.Secret, error) {
	if s.secret != nil {
		return s.secret, nil
	}
	sec := &apiv1.Secret{}
	data, err := s.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, sec)
	if err != nil {
		return nil, err
	}
	s.secret = sec
	return sec, nil
}
