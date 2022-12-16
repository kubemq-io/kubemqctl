package operator

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type roleBindingManager struct {
	*client.Client
}

func (m *roleBindingManager) CreateOrUpdateRoleBinding(roleBinding *rbac.RoleBinding) (*rbac.RoleBinding, bool, error) {
	found, err := m.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Get(context.Background(), roleBinding.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedRoleBinding, err := m.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Update(context.Background(), roleBinding, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return updatedRoleBinding, true, nil
	}

	newRoleBinding, err := m.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Create(context.Background(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return newRoleBinding, false, nil
}

func (m *roleBindingManager) DeleteRoleBinding(roleBinding *rbac.RoleBinding) error {
	found, err := m.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Get(context.Background(), roleBinding.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.RbacV1().RoleBindings(roleBinding.Namespace).Delete(context.Background(), roleBinding.Name, metav1.DeleteOptions{})
	}
	return nil
}

func (m *roleBindingManager) GetRoleBinding(name, namespace string) (*rbac.RoleBinding, error) {
	role, err := m.ClientSet.RbacV1().RoleBindings(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
