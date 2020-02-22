package operator

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type roleManager struct {
	*client.Client
}

func (m *roleManager) CreateOrUpdateRole(role *rbac.Role) (*rbac.Role, bool, error) {
	found, err := m.ClientSet.RbacV1().Roles(role.Namespace).Get(role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedRole, err := m.ClientSet.RbacV1().Roles(role.Namespace).Update(role)
		if err != nil {
			return nil, true, err
		}
		return updatedRole, true, nil
	}

	newRole, err := m.ClientSet.RbacV1().Roles(role.Namespace).Create(role)
	if err != nil {
		return nil, false, err
	}
	return newRole, false, nil
}

func (m *roleManager) DeleteRole(role *rbac.Role) error {
	found, err := m.ClientSet.RbacV1().Roles(role.Namespace).Get(role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.RbacV1().Roles(role.Namespace).Delete(role.Name, metav1.NewDeleteOptions(0))
	}
	return nil

}

func (m *roleManager) GetRole(name, namespace string) (*rbac.Role, error) {
	role, err := m.ClientSet.RbacV1().Roles(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
