package deployment

import (
	"github.com/ghodss/yaml"
	net "k8s.io/api/networking/v1beta1"
)

var defaultKubeMQSIngressTemplate = `
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.AppName}}
    deployment.id: {{.Id}}
spec:
  backend:
    serviceName: {{.ServiceName}}
    servicePort: {{.ServicePort}}
`

type IngressConfig struct {
	Id          string
	Name        string
	Namespace   string
	AppName     string
	ServiceName string
	ServicePort int
	ingress     *net.Ingress
}

func ImportIngress(spec []byte) (*IngressConfig, error) {
	ingress := &net.Ingress{}
	err := yaml.Unmarshal(spec, ingress)
	if err != nil {
		return nil, err
	}
	return &IngressConfig{
		Id:          "",
		Name:        ingress.Name,
		Namespace:   ingress.Namespace,
		AppName:     "",
		ServiceName: "",
		ServicePort: 0,
		ingress:     ingress,
	}, nil
}

func NewIngressConfig(id, name, namespace, appName string) *IngressConfig {
	return &IngressConfig{
		Id:          id,
		Name:        name,
		Namespace:   namespace,
		AppName:     appName,
		ServiceName: "",
		ServicePort: 0,
		ingress:     nil,
	}
}

func DefaultIngressConfig() map[string]*IngressConfig {
	ings := make(map[string]*IngressConfig)
	return ings
}

func (i *IngressConfig) SetServiceName(value string) *IngressConfig {
	i.ServiceName = value
	return i
}

func (i *IngressConfig) SetServicePort(value int) *IngressConfig {
	i.ServicePort = value
	return i
}

func (i *IngressConfig) Spec() ([]byte, error) {
	if i.ingress == nil {
		t := NewTemplate(defaultKubeMQSIngressTemplate, i)
		return t.Get()
	}
	return yaml.Marshal(i.ingress)
}
func (i *IngressConfig) Set(value *net.Ingress) *IngressConfig {
	i.ingress = value
	return i
}
func (i *IngressConfig) Get() (*net.Ingress, error) {
	if i.ingress != nil {
		return i.ingress, nil
	}
	ingress := &net.Ingress{}
	data, err := i.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, ingress)
	if err != nil {
		return nil, err
	}
	i.ingress = ingress
	return ingress, nil
}
