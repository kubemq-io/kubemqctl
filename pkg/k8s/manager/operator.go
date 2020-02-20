package manager

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/operator"
)

func (m *Manager) DeployOperator(namespace string) error {
	if namespace == "" {
		namespace = defaultNamespace
	}
	err := m.client.CreateIfNotPresentNamespace(namespace)
	if err != nil {
		return err
	}

	bundle, err := operator.CreateBundle(defaultOperatorName, namespace)
	if err != nil {
		return err
	}

	deployedBundle, _, err := m.client.CreateOrUpdateBundle(bundle)
	if err != nil {
		return err
	}
	return deployedBundle.Validate()
}

func (m *Manager) DeleteOperator(bundle *operator.Bundle) error {
	_, err := m.client.DeleteBundle(bundle)
	if err != nil {
		return err
	}
	return nil
}
func (m *Manager) GetOperator(namespace string) (*operator.Bundle, error) {
	if namespace == "" {
		namespace = defaultNamespace
	}

	bundle, err := m.client.GetBundle(defaultOperatorName, namespace)
	if err != nil {
		return nil, err
	}
	return bundle, nil
}

func (m *Manager) GetOperatorList() ([]string, error) {
	nsList, err := m.client.GetNamespaceList()
	if err != nil {
		return nil, err
	}
	var bundleList []string
	for _, namespace := range nsList {
		bundle, err := m.GetOperator(namespace)
		if err != nil {
			return nil, err
		}
		if err := bundle.Validate(); err == nil {
			bundleList = append(bundleList, namespace)
		}
	}

	return bundleList, nil
}
