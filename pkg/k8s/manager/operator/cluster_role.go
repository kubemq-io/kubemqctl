package operator

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterRoleManager struct {
	*client.Client
}

func (m *clusterRoleManager) CreateOrUpdateClusterRole(role *rbac.ClusterRole) (*rbac.ClusterRole, bool, error) {
	found, err := m.ClientSet.RbacV1().ClusterRoles().Get(context.Background(), role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedRole, err := m.ClientSet.RbacV1().ClusterRoles().Update(context.Background(), role, metav1.UpdateOptions{})
		if err != nil {
			return nil, true, err
		}
		return updatedRole, true, nil
	}

	newRole, err := m.ClientSet.RbacV1().ClusterRoles().Create(context.Background(), role, metav1.CreateOptions{})
	if err != nil {
		return nil, false, err
	}
	return newRole, false, nil
}

func (m *clusterRoleManager) DeleteClusterRole(role *rbac.ClusterRole) error {
	found, err := m.ClientSet.RbacV1().ClusterRoles().Get(context.Background(), role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.RbacV1().ClusterRoles().Delete(context.Background(), role.Name, metav1.DeleteOptions{})
	}
	return nil

}

func (m *clusterRoleManager) GetClusterRole(name string) (*rbac.ClusterRole, error) {
	role, err := m.ClientSet.RbacV1().ClusterRoles().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
