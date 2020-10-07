package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"sort"
)

type KubemqClusters struct {
	Items []*kubemqcluster.KubemqCluster
	m     map[string]*kubemqcluster.KubemqCluster
	list  []string
}

func newKubemqClusters(items []*kubemqcluster.KubemqCluster) *KubemqClusters {
	k := &KubemqClusters{
		Items: items,
		m:     map[string]*kubemqcluster.KubemqCluster{},
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

func (k *KubemqClusters) Cluster(name string) *kubemqcluster.KubemqCluster {
	return k.m[name]
}
func (k *KubemqClusters) List() []string {
	return k.list
}
