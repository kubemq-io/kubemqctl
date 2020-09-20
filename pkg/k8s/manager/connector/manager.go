package connector

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Manager struct {
	client *client.Client
}

func NewManager(c *client.Client) (*Manager, error) {

	return &Manager{
		client: c,
	}, nil
}

func (m *Manager) CreateOrUpdateKubemqConnector(connector *kubemqconnector.KubemqConnector) (*kubemqconnector.KubemqConnector, bool, error) {
	found, err := m.client.ClientV1Alpha1.KubemqConnector(connector.Namespace).Get(connector.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		connector.ResourceVersion = found.ResourceVersion
		updatedDashboard, err := m.client.ClientV1Alpha1.KubemqConnector(connector.Namespace).Update(connector)
		if err != nil {
			return nil, true, err
		}
		return updatedDashboard, true, nil
	}

	newDashboard, err := m.client.ClientV1Alpha1.KubemqConnector(connector.Namespace).Create(connector)
	if err != nil {
		return nil, false, err
	}
	return newDashboard, false, nil
}

func (m *Manager) DeleteKubemqConnector(connector *kubemqconnector.KubemqConnector) error {
	if connector == nil {
		return nil
	}
	return m.client.ClientV1Alpha1.KubemqConnector(connector.Namespace).Delete(connector.Name, metav1.NewDeleteOptions(0))

}

func (m *Manager) GetConnector(name, namespace string) (*kubemqconnector.KubemqConnector, error) {
	connector, err := m.client.ClientV1Alpha1.KubemqConnector(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return connector, nil
}

func (m *Manager) GetKubemqConnectors() (*KubemqConnectors, error) {
	nsList, err := m.client.GetNamespaceList()
	if err != nil {
		return nil, err
	}
	var list []*kubemqconnector.KubemqConnector
	for i := 0; i < len(nsList); i++ {
		connectors, err := m.client.ClientV1Alpha1.KubemqConnector(nsList[i]).List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(connectors.Items); i++ {
			list = append(list, &connectors.Items[i])
		}
	}
	return newKubemqConnectors(list), nil
}
