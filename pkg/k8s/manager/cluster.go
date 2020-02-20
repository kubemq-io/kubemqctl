package manager

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
)

func (m *Manager) DeployKubeMQCluster(deployment *cluster.KubemqCluster) error {
	fmt.Println(deployment.String())

	return nil
}
