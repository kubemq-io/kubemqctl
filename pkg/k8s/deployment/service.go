package deployment

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
)

var defaultKubeMQServiceTemplate = `
apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
  labels:
    app: {{.AppName}}
    deployment.id: {{.Id}}
spec:
  ports:
    - name: {{.PortName}}
      port: {{.ContainerPort}}
      protocol: TCP
      targetPort: {{.TargetPort}}
  sessionAffinity: None
  type: {{.Type}}
  selector:
    app: {{.AppName}}
`

type ServiceConfig struct {
	Id            string
	Name          string
	Namespace     string
	AppName       string
	Type          string
	ContainerPort int
	TargetPort    int
	PortName      string
	service       *apiv1.Service
}

func ImportServiceConfig(spec []byte) (*ServiceConfig, error) {
	svc := &apiv1.Service{}
	err := yaml.Unmarshal(spec, svc)
	if err != nil {
		return nil, err
	}
	return &ServiceConfig{
		Id:            "",
		Name:          svc.Name,
		Namespace:     svc.Namespace,
		AppName:       "",
		Type:          "",
		ContainerPort: 0,
		TargetPort:    0,
		PortName:      "",
		service:       svc,
	}, nil
}

func NewServiceConfig(id, name, namespace, appName string) *ServiceConfig {
	return &ServiceConfig{
		Id:            id,
		Name:          name,
		Namespace:     namespace,
		AppName:       appName,
		Type:          "",
		ContainerPort: 0,
		TargetPort:    0,
		PortName:      "",
	}
}

func DefaultServiceConfig(id, namespace, appName string) map[string]*ServiceConfig {
	list := map[string]*ServiceConfig{}
	list["grpc"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-grpc",
		Namespace:     namespace,
		AppName:       appName,
		Type:          "ClusterIP",
		ContainerPort: 50000,
		TargetPort:    50000,
		PortName:      "grpc-port",
	}
	list["rest"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-rest",
		Namespace:     namespace,
		AppName:       appName,
		Type:          "ClusterIP",
		ContainerPort: 9090,
		TargetPort:    9090,
		PortName:      "rest-port",
	}
	list["api"] = &ServiceConfig{
		Id:            id,
		Name:          appName + "-api",
		Namespace:     namespace,
		AppName:       appName,
		Type:          "ClusterIP",
		ContainerPort: 8080,
		TargetPort:    8080,
		PortName:      "rest-port",
	}
	list["internal"] = &ServiceConfig{
		Id:            id,
		Name:          appName,
		Namespace:     namespace,
		AppName:       appName,
		Type:          "ClusterIP",
		ContainerPort: 5228,
		TargetPort:    5228,
		PortName:      "cluster-port",
	}
	return list
}

func (s *ServiceConfig) SetType(value string) *ServiceConfig {
	s.Type = value
	return s
}
func (s *ServiceConfig) SetContainerPort(value int) *ServiceConfig {
	s.ContainerPort = value
	return s
}
func (s *ServiceConfig) SetTargetPort(value int) *ServiceConfig {
	s.TargetPort = value
	return s
}
func (s *ServiceConfig) SetPortName(value string) *ServiceConfig {
	s.PortName = value
	return s
}
func (s *ServiceConfig) Spec() ([]byte, error) {
	if s.service == nil {
		t := NewTemplate(defaultKubeMQServiceTemplate, s)
		return t.Get()
	}
	return yaml.Marshal(s.service)
}
func (s *ServiceConfig) Set(value *apiv1.Service) *ServiceConfig {
	s.service = value
	return s
}
func (s *ServiceConfig) Get() (*apiv1.Service, error) {
	if s.service != nil {
		return s.service, nil
	}
	svc := &apiv1.Service{}
	data, err := s.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, svc)
	if err != nil {
		return nil, err
	}
	s.service = svc
	return svc, nil
}
