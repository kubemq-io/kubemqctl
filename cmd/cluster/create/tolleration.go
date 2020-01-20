package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
)

type deployTolerationOptions struct {
	enabled bool
}

func defaultTolerationOptions(cmd *cobra.Command) *deployTolerationOptions {
	o := &deployTolerationOptions{
		enabled: false,
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "toleration-enabled", "", false, "enable statefulset toleration  configuration")
	return o
}

func (o *deployTolerationOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployTolerationOptions) complete() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployTolerationOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployTolerationOptions {
	//TODO - Add code
	return o
}
