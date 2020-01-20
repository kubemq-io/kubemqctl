package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
)

type deployNodeSelectorOptions struct {
	enabled bool
}

func defaultNodeSelectorOptions(cmd *cobra.Command) *deployNodeSelectorOptions {
	o := &deployNodeSelectorOptions{
		enabled: false,
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "node-selectors-enabled", "", false, "enable statefulset node selectors configuration")
	return o
}

func (o *deployNodeSelectorOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployNodeSelectorOptions) complete() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployNodeSelectorOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployNodeSelectorOptions {
	//TODO - Add code
	return o
}
