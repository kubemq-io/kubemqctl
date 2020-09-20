package operator

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterRoleManager struct {
	*client.Client
}

func (m *clusterRoleManager) CreateOrUpdateClusterRole(role *rbac.ClusterRole) (*rbac.ClusterRole, bool, error) {
	found, err := m.ClientSet.RbacV1().ClusterRoles().Get(role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		updatedRole, err := m.ClientSet.RbacV1().ClusterRoles().Update(role)
		if err != nil {
			return nil, true, err
		}
		return updatedRole, true, nil
	}

	newRole, err := m.ClientSet.RbacV1().ClusterRoles().Create(role)
	if err != nil {
		return nil, false, err
	}
	return newRole, false, nil
}

func (m *clusterRoleManager) DeleteClusterRole(role *rbac.ClusterRole) error {
	found, err := m.ClientSet.RbacV1().ClusterRoles().Get(role.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		return m.ClientSet.RbacV1().ClusterRoles().Delete(role.Name, metav1.NewDeleteOptions(0))
	}
	return nil

}

func (m *clusterRoleManager) GetClusterRole(name string) (*rbac.ClusterRole, error) {
	role, err := m.ClientSet.RbacV1().ClusterRoles().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return role, nil
}
