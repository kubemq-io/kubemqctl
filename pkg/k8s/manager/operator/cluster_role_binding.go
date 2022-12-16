package operator

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterRoleBindingManager struct {
	*client.Client
}

func (m *clusterRoleBindingManager) CreateOrUpdateClusterRoleBinding(roleBinding *rbac.ClusterRoleBinding) (*rbac.ClusterRoleBinding, bool, error) {
	found, err := m.ClientSet.RbacV1().ClusterRoleBindings().Get(context.Background(), roleBinding.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedRoleBinding, err := m.ClientSet.RbacV1().ClusterRoleBindings().Update(context.Background(), roleBinding, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return updatedRoleBinding, true, nil
	}

	newRoleBinding, err := m.ClientSet.RbacV1().ClusterRoleBindings().Create(context.Background(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return newRoleBinding, false, nil
}

func (m *clusterRoleBindingManager) DeleteClusterRoleBinding(roleBinding *rbac.ClusterRoleBinding) error {
	found, err := m.ClientSet.RbacV1().ClusterRoleBindings().Get(context.Background(), roleBinding.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.RbacV1().ClusterRoleBindings().Delete(context.Background(), roleBinding.Name, metav1.DeleteOptions{})
	}
	return nil
}

func (m *clusterRoleBindingManager) GetClusterRoleBinding(name, namespace string) (*rbac.ClusterRoleBinding, error) {
	role, err := m.ClientSet.RbacV1().ClusterRoleBindings().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
