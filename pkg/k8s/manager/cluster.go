package manager

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
)

func (m *Manager) DeployKubeMQCluster(deployment *cluster.KubemqCluster) error {
	deployment.Kind = "KubemqCluster"
	deployment.APIVersion = "core.k8s.kubemq.io/v1alpha1"
	create, err := m.client.ClientV1Alpha1.KubemqClusters("kubemq").Create(deployment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(create.String())

	return nil
}
