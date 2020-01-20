package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/spf13/cobra"
)

type deployResourceOptions struct {
	enabled        bool
	limitsCpu      string
	limitsMemory   string
	requestsCpu    string
	requestsMemory string
}

func defaultResourceOptions(cmd *cobra.Command) *deployResourceOptions {
	o := &deployResourceOptions{
		enabled:        false,
		limitsCpu:      "",
		limitsMemory:   "",
		requestsCpu:    "",
		requestsMemory: "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "resources-enabled", "", false, "enable resources configuration")
	cmd.PersistentFlags().StringVarP(&o.limitsCpu, "resources-limits-key-cpu", "", "", "set resources limits cpu ")
	cmd.PersistentFlags().StringVarP(&o.limitsMemory, "resources-limits-key-memory", "", "", "set resources limits memory")
	cmd.PersistentFlags().StringVarP(&o.requestsCpu, "resources-requests-key-cpu", "", "", "set resources requests cpu")
	cmd.PersistentFlags().StringVarP(&o.requestsMemory, "resources-requests-memory", "", "", "set resources request memory")

	return o
}

func (o *deployResourceOptions) validate() error {
	if !o.enabled {
		return nil
	}
	if o.limitsCpu == "" {
		return fmt.Errorf("error setting resources configuration, missing limits cpu data")
	}
	if o.limitsMemory == "" {
		return fmt.Errorf("error setting resources configuration, missing limits memory data")
	}
	if o.requestsCpu == "" {
		return fmt.Errorf("error setting resources configuration, missing requests cpu data")
	}
	if o.requestsMemory == "" {
		return fmt.Errorf("error setting resources configuration, missing requests memory data")
	}
	return nil
}

func (o *deployResourceOptions) complete() error {
	return nil
}

func (o *deployResourceOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployResourceOptions {
	if !o.enabled {
		return o
	}
	//TODO - add code

	return o
}
