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

func newKubemqConnectors(items []*kubemqconnector.KubemqConnector) *KubemqConnectors {
	k := &KubemqConnectors{
		items: items,
		m:     map[string]*kubemqconnector.KubemqConnector{},
		list:  []string{},
	}
	for i := 0; i < len(items); i++ {
		pair := fmt.Sprintf("%s/%s", items[i].Namespace, items[i].Name)
		k.list = append(k.list, pair)
		k.m[pair] = items[i]
	}
	sort.Strings(k.list)
	return k
}

func (k *KubemqConnectors) Connector(name string) *kubemqconnector.KubemqConnector {
	return k.m[name]
}
func (k *KubemqConnectors) List() []string {
	return k.list
}
