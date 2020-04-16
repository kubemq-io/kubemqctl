package dashboard

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
	"sort"
)

type KubemqDashboards struct {
	items []*kubemqdashboard.KubemqDashboard
	m     map[string]*kubemqdashboard.KubemqDashboard
	list  []string
}

func newKubemqDashboards(items []*kubemqdashboard.KubemqDashboard) *KubemqDashboards {
	k := &KubemqDashboards{
		items: items,
		m:     map[string]*kubemqdashboard.KubemqDashboard{},
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

func (k *KubemqDashboards) Dashboard(name string) *kubemqdashboard.KubemqDashboard {
	return k.m[name]
}
func (k *KubemqDashboards) List() []string {
	return k.list
}
