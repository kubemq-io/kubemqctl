package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

type deployNodeSelectorOptions struct {
	keys map[string]string
}

func defaultNodeSelectorOptions(cmd *cobra.Command) *deployNodeSelectorOptions {
	o := &deployNodeSelectorOptions{
		keys: map[string]string{},
	}
	cmd.PersistentFlags().StringToStringVarP(&o.keys, "node-selectors-keys", "", map[string]string{}, "set statefulset node selectors key-value (map)")
	return o
}

func (o *deployNodeSelectorOptions) validate() error {

	return nil
}

func (o *deployNodeSelectorOptions) complete() error {
	return nil
}

func (o *deployNodeSelectorOptions) setConfig(deployment *cluster.KubemqCluster) *deployNodeSelectorOptions {
	if len(o.keys) > 0 {
		return nil
	}
	deployment.Spec.NodeSelectors = &cluster.NodeSelectorConfig{Keys: o.keys}

	return o
}
