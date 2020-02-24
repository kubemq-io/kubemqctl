package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

type deployNodeSelectorOptions struct {
	keys map[string]string
}

func setNodeSelectorOptions(cmd *cobra.Command) *deployNodeSelectorOptions {
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

func (o *deployNodeSelectorOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployNodeSelectorOptions {
	if len(o.keys) == 0 {
		return nil
	}
	deployment.Spec.NodeSelectors = &kubemqcluster.NodeSelectorConfig{Keys: o.keys}

	return o
}
