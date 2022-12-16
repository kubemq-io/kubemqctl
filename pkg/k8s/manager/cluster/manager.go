package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
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

func (m *Manager) CreateOrUpdateKubemqCluster(cluster *kubemqcluster.KubemqCluster) (*kubemqcluster.KubemqCluster, bool, error) {
	found, err := m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).Get(cluster.Name, metav1.GetOptions{})
	if err == nil && found != nil {
		cluster.ResourceVersion = found.ResourceVersion
		updatedCluster, err := m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).Update(cluster)
		if err != nil {
			return nil, true, err
		}
		return updatedCluster, true, nil
	}

	newCluster, err := m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).Create(cluster)
	if err != nil {
		return nil, false, err
	}
	return newCluster, false, nil
}

func (m *Manager) DeleteKubemqCluster(cluster *kubemqcluster.KubemqCluster) error {
	if cluster == nil {
		return nil
	}
	return m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).Delete(cluster.Name, metav1.NewDeleteOptions(0))
}

func (m *Manager) GetCluster(name, namespace string) (*kubemqcluster.KubemqCluster, error) {
	cluster, err := m.client.ClientV1Beta1.KubemqClusters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cluster, nil
}

func (m *Manager) GetKubemqClusters() (*KubemqClusters, error) {
	var list []*kubemqcluster.KubemqCluster
	values, err := m.client.ClientV1Beta1.KubemqClusters("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(values.Items); i++ {
		list = append(list, &values.Items[i])
	}
	connectors := newKubemqClusters().SetItems(list)
	return connectors, nil
}

func (m *Manager) ScaleKubemqCluster(cluster *kubemqcluster.KubemqCluster, scaleTo int32) error {
	if cluster == nil {
		return nil
	}
	scale, err := m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).GetScale(cluster.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	scale.Spec.Replicas = scaleTo
	_, err = m.client.ClientV1Beta1.KubemqClusters(cluster.Namespace).UpdateScale(cluster.Name, scale)
	return err
}
