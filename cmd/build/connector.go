package build

import (
	"fmt"
	connectorbuilder "github.com/kubemq-hub/builder/connector"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConnectorsBuilder struct {
	deployments    []*kubemqconnector.KubemqConnector
	resources      *resources
	defaultOptions common.DefaultOptions
}

func newConnectorsBuilder() *ConnectorsBuilder {
	return &ConnectorsBuilder{}
}
func (c *ConnectorsBuilder) SetResources(value *resources) *ConnectorsBuilder {
	c.resources = value
	return c
}

func (c *ConnectorsBuilder) render() error {
	connector, err := connectorbuilder.NewConnector(nil).
		SetDefaultOptions(c.defaultOptions).
		Render()
	if err != nil {
		return err
	}
	deployment := &kubemqconnector.KubemqConnector{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqConnector",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s", connector.Name, connector.Type),
			Namespace: connector.Namespace,
		},
		Spec: kubemqconnector.KubemqConnectorSpec{
			Replicas:    new(int32),
			Type:        connector.Type,
			Image:       connector.Image,
			Config:      connector.Config,
			NodePort:    int32(connector.NodePort),
			ServiceType: connector.ServiceType,
		},
		Status: kubemqconnector.KubemqConnectorStatus{},
	}
	c.deployments = append(c.deployments, deployment)
	return nil
}
func (c *ConnectorsBuilder) add() error {
	utils.Println("Adding new KubeMQ Connector:")
	if err := c.render(); err != nil {
		return err
	}
	utils.Println("Adding new KubeMQ Connector completed successfully!")
	return nil
}
