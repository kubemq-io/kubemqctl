package operator

import (
	"fmt"
	"sort"

	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
)

type Operators struct {
	items []*operator.Deployment
	m     map[string]*operator.Deployment
	list  []string
}

func newOperators(items []*operator.Deployment) *Operators {
	k := &Operators{
		items: items,
		m:     map[string]*operator.Deployment{},
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

func (o *Operators) Deployment(name string) *operator.Deployment {
	return o.m[name]
}

func (o *Operators) List() []string {
	return o.list
}
