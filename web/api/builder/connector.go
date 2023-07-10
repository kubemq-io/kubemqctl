package builder

import (
	"fmt"
	"github.com/kubemq-io/k8s/api/v1beta1"
)

type Connector struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Metadata   struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec *v1beta1.KubemqConnectorSpec `json:"spec"`
}

func NewConnector() *Connector {
	c := &Connector{
		Kind:       "KubemqConnector",
		ApiVersion: "core.k8s.kubemq.io/v1beta1",
		Spec: &v1beta1.KubemqConnectorSpec{
			Replicas:    new(int32),
			Type:        "",
			Image:       "",
			Config:      "",
			NodePort:    0,
			ServiceType: "",
		},
	}
	*c.Spec.Replicas = 1
	return c
}

func (c *Connector) SetName(value string) *Connector {
	c.Metadata.Name = value
	return c
}

func (c *Connector) SetNamespace(value string) *Connector {
	c.Metadata.Namespace = value
	return c
}

func (c *Connector) SetType(value string) *Connector {
	switch value {
	case "bridges", "sources", "targets":
		c.Spec.Type = value
	}
	return c
}
func (c *Connector) SetServiceType(value string) *Connector {
	switch value {
	case "ClusterIP", "LoadBalancer", "NodePort":
		c.Spec.ServiceType = value
	}
	return c
}
func (c *Connector) SetConfig(value string) *Connector {
	c.Spec.Config = value
	return c
}
func (c *Connector) SetNodePort(value int32) *Connector {
	c.Spec.NodePort = value
	return c
}

func (c *Connector) Validate() error {
	if c.Metadata.Name == "" {
		return fmt.Errorf("invalid connector name")
	}
	if c.Metadata.Namespace == "" {
		return fmt.Errorf("invalid connector namespace")
	}
	if c.Spec.Type == "" {
		return fmt.Errorf("invalid connector type")
	}
	if c.Spec.ServiceType == "" {
		return fmt.Errorf("invalid connector service type")
	}
	if c.Spec.Config == "" {
		return fmt.Errorf("invalid connector config")
	}
	return nil
}
