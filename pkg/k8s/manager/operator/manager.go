package operator

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manager struct {
	*client.Client
	*crdManager
	*operatorManager
	*clusterRoleManager
	*clusterRoleBindingManager
	*serviceAccountManager
	*roleManager
	*roleBindingManager
}

func NewManager(c *client.Client) (*Manager, error) {
	return &Manager{
		Client:                    c,
		crdManager:                &crdManager{c},
		operatorManager:           &operatorManager{c},
		clusterRoleManager:        &clusterRoleManager{c},
		clusterRoleBindingManager: &clusterRoleBindingManager{c},
		serviceAccountManager:     &serviceAccountManager{c},
		roleManager:               &roleManager{c},
		roleBindingManager:        &roleBindingManager{c},
	}, nil
}
func (m *Manager) checkAndCreateNamespace(ns string) error {
	newNs := &apiv1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
		Spec:   apiv1.NamespaceSpec{},
		Status: apiv1.NamespaceStatus{},
	}
	_, _, err := m.Client.CheckAndCreateNamespace(newNs)
	return err
}

func (m *Manager) CreateOrUpdateKubemqOperator(operatorDeployment *operator.Deployment) (*operator.Deployment, bool, error) {
	err := m.checkAndCreateNamespace(operatorDeployment.Namespace)
	if err != nil {
		return nil, false, err
	}

	newBundle := &operator.Deployment{
		Name:                   operatorDeployment.Name,
		Namespace:              operatorDeployment.Namespace,
		CRDs:                   nil,
		Deployment:             nil,
		ClusterRole:            nil,
		ClusterRoleBinding:     nil,
		Role:                   nil,
		RoleBinding:            nil,
		OperatorServiceAccount: nil,
		ClusterServiceAccount:  nil,
	}
	isBundleUpdated := false
	operatorServiceAccount, isUpdated, err := m.CreateOrUpdateServiceAccount(operatorDeployment.OperatorServiceAccount)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update service account, error: %s", err.Error())
	}
	newBundle.OperatorServiceAccount = operatorServiceAccount
	if isUpdated {
		isBundleUpdated = true
	}

	clusterServiceAccount, isUpdated, err := m.CreateOrUpdateServiceAccount(operatorDeployment.ClusterServiceAccount)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update service account, error: %s", err.Error())
	}
	newBundle.ClusterServiceAccount = clusterServiceAccount
	if isUpdated {
		isBundleUpdated = true
	}

	clusterRole, isUpdated, err := m.CreateOrUpdateClusterRole(operatorDeployment.ClusterRole)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update cluster role, error: %s", err.Error())
	}
	newBundle.ClusterRole = clusterRole
	if isUpdated {
		isBundleUpdated = true
	}

	clusterRoleBinding, isUpdated, err := m.CreateOrUpdateClusterRoleBinding(operatorDeployment.ClusterRoleBinding)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update cluster role binding, error: %s", err.Error())
	}
	newBundle.ClusterRoleBinding = clusterRoleBinding
	if isUpdated {
		isBundleUpdated = true
	}

	role, isUpdated, err := m.CreateOrUpdateRole(operatorDeployment.Role)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update role, error: %s", err.Error())
	}
	newBundle.Role = role
	if isUpdated {
		isBundleUpdated = true
	}

	roleBinding, isUpdated, err := m.CreateOrUpdateRoleBinding(operatorDeployment.RoleBinding)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update role binding, error: %s", err.Error())
	}
	newBundle.RoleBinding = roleBinding
	if isUpdated {
		isBundleUpdated = true
	}

	for _, crd := range operatorDeployment.CRDs {
		newCrd, isUpdated, err := m.CreateOrUpdateCRD(crd)
		if err != nil {
			return nil, false, fmt.Errorf("error create or update crd, error: %s", err.Error())
		}
		newBundle.CRDs = append(newBundle.CRDs, newCrd)
		if isUpdated {
			isBundleUpdated = true
		}
	}

	deployment, isUpdated, err := m.CreateOrUpdateOperator(operatorDeployment.Deployment)
	if err != nil {
		return nil, false, fmt.Errorf("error create or update operator deployment, error: %s", err.Error())
	}
	newBundle.Deployment = deployment
	if isUpdated {
		isBundleUpdated = true
	}
	return newBundle, isBundleUpdated, nil
}

func (m *Manager) DeleteKubemqOperator(deployment *operator.Deployment, isAll bool) error {
	err := m.DeleteOperator(deployment.Deployment)
	if err != nil {
		return fmt.Errorf("delete operator failed, error: %s", err.Error())
	}

	if isAll {
		err = m.DeleteClusterRole(deployment.ClusterRole)
		if err != nil {
			return fmt.Errorf("delete cluster role failed, error: %s", err.Error())
		}

		err = m.DeleteClusterRoleBinding(deployment.ClusterRoleBinding)
		if err != nil {
			return fmt.Errorf("delete cluster role binding failed, error: %s", err.Error())
		}
		if deployment.Role != nil {
			err = m.DeleteRole(deployment.Role)
			if err != nil {
				return fmt.Errorf("delete role failed, error: %s", err.Error())
			}
		}

		if deployment.RoleBinding != nil {
			err = m.DeleteRoleBinding(deployment.RoleBinding)
			if err != nil {
				return fmt.Errorf("delete role binding failed, error: %s", err.Error())
			}
		}
		err = m.DeleteServiceAccount(deployment.OperatorServiceAccount)
		if err != nil {
			return fmt.Errorf("delete service acount failed, error: %s", err.Error())
		}
		err = m.DeleteServiceAccount(deployment.ClusterServiceAccount)
		if err != nil {
			return fmt.Errorf("delete service acount failed, error: %s", err.Error())
		}
	}

	return nil
}
func (m *Manager) IsKubemqOperatorExists(namespace string) bool {
	_, err := m.GetOperator("kubemq-operator", namespace)
	return err == nil
}

func (m *Manager) GetKubemqOperator(name, namespace string) (*operator.Deployment, error) {
	bundle := &operator.Deployment{
		Name:                   name,
		Namespace:              namespace,
		CRDs:                   nil,
		Deployment:             nil,
		ClusterRole:            nil,
		ClusterRoleBinding:     nil,
		Role:                   nil,
		RoleBinding:            nil,
		OperatorServiceAccount: nil,
		ClusterServiceAccount:  nil,
	}

	bundle.Deployment, _ = m.GetOperator(name, namespace)
	crd, _ := m.GetCrd("kubemqclusters.core.k8s.kubemq.io", namespace)
	if crd != nil {
		bundle.CRDs = append(bundle.CRDs, crd)
	}
	crd, _ = m.GetCrd("kubemqdashboards.core.k8s.kubemq.io", namespace)
	if crd != nil {
		bundle.CRDs = append(bundle.CRDs, crd)
	}
	crd, _ = m.GetCrd("kubemqconnectors.core.k8s.kubemq.io", namespace)
	if crd != nil {
		bundle.CRDs = append(bundle.CRDs, crd)
	}

	bundle.ClusterRole, _ = m.GetClusterRole("kubemq-operator")
	bundle.ClusterRoleBinding, _ = m.GetClusterRoleBinding(fmt.Sprintf("kubemq-operator-%s-crb", namespace), namespace)
	bundle.Role, _ = m.GetRole("kubemq-cluster", namespace)
	bundle.RoleBinding, _ = m.GetRoleBinding("kubemq-cluster", namespace)
	bundle.OperatorServiceAccount, _ = m.GetServiceAccount("kubemq-operator", namespace)
	bundle.ClusterServiceAccount, _ = m.GetServiceAccount("kubemq-cluster", namespace)
	return bundle, nil
}
func (m *Manager) GetKubemqOperators() (*Operators, error) {
	nsList, err := m.GetNamespaceList()
	if err != nil {
		return nil, err
	}
	var list []*operator.Deployment
	for i := 0; i < len(nsList); i++ {
		op, err := m.GetKubemqOperator("kubemq-operator", nsList[i])
		if err != nil {
			return nil, err
		}
		if err := op.IsValid(); err == nil {

			list = append(list, op)
		}
	}
	return newOperators(list), nil
}
