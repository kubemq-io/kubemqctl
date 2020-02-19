package manager

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
)

const (
	defaultOperatorName = "kubemq-operator"
	defaultNamespace    = "kubemq"
)

type Manager struct {
	client *client.Client
}

func NewManager(kubeConfigPath string) (*Manager, error) {
	c, err := client.NewClient(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return &Manager{
		client: c,
	}, nil
}
