package deployment

import (
	"github.com/ghodss/yaml"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var defaultNamespaceTemplate = `
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Name}}
  labels:
    deployment.id: {{.Id}}
`

type NamespaceConfig struct {
	Id        string
	Name      string
	namespace *apiv1.Namespace
}

func ImportNamespaceConfig(spec []byte) (*NamespaceConfig, error) {
	ns := &apiv1.Namespace{}
	err := yaml.Unmarshal(spec, ns)
	if err != nil {
		return nil, err
	}
	return &NamespaceConfig{
		Id:        "",
		Name:      ns.Name,
		namespace: ns,
	}, nil
}

func NewNamespaceConfig(id, name string) *NamespaceConfig {
	return &NamespaceConfig{
		Id:   id,
		Name: name,
	}
}
func DefaultNamespaceConfig(id, name string) *NamespaceConfig {
	return &NamespaceConfig{
		Id:   id,
		Name: name,
	}
}
func (n *NamespaceConfig) Spec() ([]byte, error) {
	if n.namespace == nil {
		t := NewTemplate(defaultNamespaceTemplate, n)
		return t.Get()
	}
	return yaml.Marshal(n.namespace)
}
func (n *NamespaceConfig) Set(value *apiv1.Namespace) *NamespaceConfig {
	n.namespace = value
	return n
}

func (n *NamespaceConfig) Get() (*apiv1.Namespace, error) {
	if n.namespace != nil {
		return n.namespace, nil
	}
	ns := &apiv1.Namespace{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Spec:       apiv1.NamespaceSpec{},
		Status:     apiv1.NamespaceStatus{},
	}
	data, err := n.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, ns)
	if err != nil {
		return nil, err
	}
	n.namespace = ns
	return ns, nil
}
