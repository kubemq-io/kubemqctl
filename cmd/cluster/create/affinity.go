package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
)

type deployAffinityOptions struct {
	enabled bool
}

func defaultAffinityOptions(cmd *cobra.Command) *deployAffinityOptions {
	o := &deployAffinityOptions{
		enabled: false,
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "toleration-enabled", "", false, "enable statefulset affinity configuration")
	return o
}

func (o *deployAffinityOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployAffinityOptions) complete() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployAffinityOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployAffinityOptions {
	//TODO - Add code
	return o
}
