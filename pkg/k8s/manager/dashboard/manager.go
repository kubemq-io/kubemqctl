package dashboard

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
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

func (m *Manager) CreateOrUpdateKubemqDashboard(dashboard *kubemqdashboard.KubemqDashboard) (*kubemqdashboard.KubemqDashboard, bool, error) {
	found, err := m.client.ClientV1Alpha1.KubemqDashboard(dashboard.Namespace).Get(dashboard.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		dashboard.ResourceVersion = found.ResourceVersion
		updatedDashboard, err := m.client.ClientV1Alpha1.KubemqDashboard(dashboard.Namespace).Update(dashboard)
		if err != nil {
			return nil, true, err
		}
		return updatedDashboard, true, nil
	}

	newDashboard, err := m.client.ClientV1Alpha1.KubemqDashboard(dashboard.Namespace).Create(dashboard)
	if err != nil {
		return nil, false, err
	}
	return newDashboard, false, nil
}

func (m *Manager) DeleteKubemqDashboard(dashboard *kubemqdashboard.KubemqDashboard) error {
	if dashboard == nil {
		return nil
	}
	return m.client.ClientV1Alpha1.KubemqDashboard(dashboard.Namespace).Delete(dashboard.Name, metav1.NewDeleteOptions(0))

}

func (m *Manager) GetDashboard(name, namespace string) (*kubemqdashboard.KubemqDashboard, error) {
	dashboard, err := m.client.ClientV1Alpha1.KubemqDashboard(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}

func (m *Manager) GetKubemqDashboardes() (*KubemqDashboards, error) {
	nsList, err := m.client.GetNamespaceList()
	if err != nil {
		return nil, err
	}
	var list []*kubemqdashboard.KubemqDashboard
	for i := 0; i < len(nsList); i++ {
		dashboards, err := m.client.ClientV1Alpha1.KubemqDashboard(nsList[i]).List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(dashboards.Items); i++ {
			list = append(list, &dashboards.Items[i])
		}
	}
	return newKubemqDashboards(list), nil
}
