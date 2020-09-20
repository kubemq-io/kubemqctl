package operator

import (
	"github.com/ghodss/yaml"
	rbac "k8s.io/api/rbac/v1"
)

var clusterRoleBinding = `
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
subjects:
- kind: ServiceAccount
  name: kubemq-operator
  namespace: {{.Namespace}}
roleRef:
  kind: ClusterRole
  name: kubemq-operator
  apiGroup: rbac.authorization.k8s.io
`

type ClusterRoleBinding struct {
	Name      string
	Namespace string
	role      *rbac.ClusterRoleBinding
}

func CreateClusterRoleBinding(name, namespace string) *ClusterRoleBinding {
	return &ClusterRoleBinding{
		Name:      name,
		Namespace: namespace,
		role:      nil,
	}
}
func (rb *ClusterRoleBinding) Spec() ([]byte, error) {
	t := NewTemplate(clusterRoleBinding, rb)
	return t.Get()
}
func (rb *ClusterRoleBinding) Get() (*rbac.ClusterRoleBinding, error) {
	if rb.role != nil {
		return rb.role, nil
	}
	roleBinding := &rbac.ClusterRoleBinding{}
	data, err := rb.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, roleBinding)
	if err != nil {
		return nil, err
	}
	rb.role = roleBinding
	return roleBinding, nil
}
