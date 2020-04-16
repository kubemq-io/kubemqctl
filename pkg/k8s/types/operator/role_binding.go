package operator

import (
	"github.com/ghodss/yaml"
	rbac "k8s.io/api/rbac/v1"
)

var roleBinding = `
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
subjects:
- kind: ServiceAccount
  name: {{.Name}}
roleRef:
  kind: Role
  name: {{.Name}}
  apiGroup: rbac.authorization.k8s.io
`

type RoleBinding struct {
	Name      string
	Namespace string
	role      *rbac.RoleBinding
}

func CreateRoleBinding(name, namespace string) *RoleBinding {
	return &RoleBinding{
		Name:      name,
		Namespace: namespace,
		role:      nil,
	}
}
func (rb *RoleBinding) Spec() ([]byte, error) {
	t := NewTemplate(roleBinding, rb)
	return t.Get()
}
func (rb *RoleBinding) Get() (*rbac.RoleBinding, error) {
	if rb.role != nil {
		return rb.role, nil
	}
	roleBinding := &rbac.RoleBinding{}
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
