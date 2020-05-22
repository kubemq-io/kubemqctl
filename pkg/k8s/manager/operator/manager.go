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
	*roleManager
	*roleBindingManager
	*serviceAccountManager
}

func NewManager(c *client.Client) (*Manager, error) {
	return &Manager{
		Client:                c,
		crdManager:            &crdManager{c},
		operatorManager:       &operatorManager{c},
		roleManager:           &roleManager{c},
		roleBindingManager:    &roleBindingManager{c},
		serviceAccountManager: &serviceAccountManager{c},
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
		for _, crd := range deployment.CRDs {
			err := m.DeleteCrd(crd)
			if err != nil {
				return fmt.Errorf("delete crd failed, error: %s", err.Error())
			}
		}

		err = m.DeleteRole(deployment.Role)
		if err != nil {
			return fmt.Errorf("delete role failed, error: %s", err.Error())
		}

		err = m.DeleteRoleBinding(deployment.RoleBinding)
		if err != nil {
			return fmt.Errorf("delete role binding failed, error: %s", err.Error())
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

func (m *Manager) GetKubemqOperator(name, namespace string) (*operator.Deployment, error) {
	bundle := &operator.Deployment{
		Name:                   name,
		Namespace:              namespace,
		CRDs:                   nil,
		Deployment:             nil,
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
	bundle.Role, _ = m.GetRole(name, namespace)
	bundle.RoleBinding, _ = m.GetRoleBinding(name, namespace)
	bundle.OperatorServiceAccount, _ = m.GetServiceAccount(name, namespace)
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
