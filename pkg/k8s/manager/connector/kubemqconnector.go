package connector

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqconnector"
	"sort"
)

type KubemqConnectors struct {
	items []*kubemqconnector.KubemqConnector
	m     map[string]*kubemqconnector.KubemqConnector
	list  []string
}

func newKubemqConnectors() *KubemqConnectors {
	return &KubemqConnectors{
		items: []*kubemqconnector.KubemqConnector{},
		m:     map[string]*kubemqconnector.KubemqConnector{},
		list:  []string{},
	}
}

func (c *KubemqConnectors) SetItems(items []*kubemqconnector.KubemqConnector) *KubemqConnectors {
	c.items = items
	for i := 0; i < len(items); i++ {
		pair := fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name)
		c.list = append(c.list, pair)
		c.m[pair] = items[i]
	}
	sort.Strings(c.list)
	return c
}
func (c *KubemqConnectors) Connector(name string) *kubemqconnector.KubemqConnector {
	return c.m[name]
}
func (c *KubemqConnectors) Items() []*kubemqconnector.KubemqConnector {
	return c.items
}

func (c *KubemqConnectors) List() []string {
	return c.list
}
