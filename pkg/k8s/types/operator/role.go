package operator

import (
	"github.com/ghodss/yaml"
	rbac "k8s.io/api/rbac/v1"
)

var role = `
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
rules:
- apiGroups:
  - ""
  resources:
  - pods
  - services
  - services/finalizers
  - endpoints
  - persistentvolumeclaims
  - events
  - configmaps
  - secrets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps
  resources:
  - deployments
  - daemonsets
  - replicasets
  - statefulsets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - monitoring.coreos.com
  resources:
  - servicemonitors
  verbs:
  - get
  - create
- apiGroups:
  - apps
  resourceNames:
  - kubemq-operator
  resources:
  - deployments/finalizers
  verbs:
  - update
- apiGroups:
  - ""
  resources:
  - pods
  verbs:
  - get
- apiGroups:
  - apps
  resources:
  - replicasets
  - deployments
  verbs:
  - get
- apiGroups:
  - core.k8s.kubemq.io
  resources:
  - '*'
  - kubemqclusters
  - kubemqdashboards
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
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
