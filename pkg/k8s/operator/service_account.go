package operator

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
)

var serviceAccount = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
`

type ServiceAccount struct {
	Name           string
	Namespace      string
	serviceAccount *apiv1.ServiceAccount
}

func CreateServiceAccount(name, namespace string) *ServiceAccount {
	return &ServiceAccount{
		Name:           name,
		Namespace:      namespace,
		serviceAccount: nil,
	}
}
func (sa *ServiceAccount) Spec() ([]byte, error) {
	t := NewTemplate(serviceAccount, sa)
	return t.Get()
}
func (sa *ServiceAccount) Get() (*apiv1.ServiceAccount, error) {
	if sa.serviceAccount != nil {
		return sa.serviceAccount, nil
	}
	serviceAccount := &apiv1.ServiceAccount{}
	data, err := sa.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, serviceAccount)
	if err != nil {
		return nil, err
	}
	sa.serviceAccount = serviceAccount
	return serviceAccount, nil
}
