package operator

import (
	"github.com/ghodss/yaml"
	rbac "k8s.io/api/rbac/v1"
)

var clusterRole = `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubemq-operator
rules:
  - apiGroups:
      - ""
    resources:
      - pods
      - services
      - endpoints
      - persistentvolumeclaims
      - events
      - configmaps
      - serviceaccounts
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
      - apiextensions.k8s.io
    resources:
      - customresourcedefinitions
    verbs:
      - patch
      - update
      - create
      - get
  - apiGroups:
      - apps
    resources:
      - deployments
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
      - core.k8s.kubemq.io
    resources:
      - "*"
    verbs:
      - create
      - delete
      - get
      - list
      - patch
      - update
      - watch
`

type ClusterRole struct {
	Name string
	role *rbac.ClusterRole
}

func CreateClusterRole(name string) *ClusterRole {
	return &ClusterRole{
		Name: name,
		role: nil,
	}
}
func (rb *ClusterRole) Spec() ([]byte, error) {
	t := NewTemplate(clusterRole, rb)
	return t.Get()
}
func (rb *ClusterRole) Get() (*rbac.ClusterRole, error) {
	if rb.role != nil {
		return rb.role, nil
	}
	role := &rbac.ClusterRole{}
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
