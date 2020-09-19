package operator

import (
	"github.com/ghodss/yaml"
	rbac "k8s.io/api/rbac/v1"
)

var role = `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kubemq-cluster
  namespace: {{.Namespace}}
rules:
  - apiGroups:
      - security.openshift.io
    resources:
      - securitycontextconstraints
    verbs:
      - use
      - delete
      - get
      - list
      - patch
      - update
      - watch
    resourceNames:
      - privileged
`

type Role struct {
	Name      string
	Namespace string
	role      *rbac.Role
}

func CreateRole(name, namespace string) *Role {
	return &Role{
		Name:      name,
		Namespace: namespace,
		role:      nil,
	}
}
func (rb *Role) Spec() ([]byte, error) {
	t := NewTemplate(role, rb)
	return t.Get()
}
func (rb *Role) Get() (*rbac.Role, error) {
	if rb.role != nil {
		return rb.role, nil
	}
	role := &rbac.Role{}
	data, err := rb.Spec()
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, role)
	if err != nil {
		return nil, err
	}
	rb.role = role
	return role, nil
}
