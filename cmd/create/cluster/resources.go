package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

type deployResourceOptions struct {
	enabled                  bool
	limitsCpu                string
	limitsMemory             string
	limitsEphemeralStorage   string
	requestsCpu              string
	requestsMemory           string
	requestsEphemeralStorage string
}

func setResourceOptions(cmd *cobra.Command) *deployResourceOptions {
	o := &deployResourceOptions{
		enabled:                  false,
		limitsCpu:                "",
		limitsMemory:             "",
		limitsEphemeralStorage:   "",
		requestsCpu:              "",
		requestsMemory:           "",
		requestsEphemeralStorage: "",
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "resources-enabled", "", false, "enable resources configuration")
	cmd.PersistentFlags().StringVarP(&o.limitsCpu, "resources-limits-cpu", "", "2", "set resources limits cpu ")
	cmd.PersistentFlags().StringVarP(&o.limitsMemory, "resources-limits-memory", "", "2Gi", "set resources limits memory")
	cmd.PersistentFlags().StringVarP(&o.limitsEphemeralStorage, "resources-limits-ephemeral-storage", "", "", "set resources limits ephemeral-storage")
	cmd.PersistentFlags().StringVarP(&o.requestsCpu, "resources-requests-cpu", "", "2", "set resources requests cpu")
	cmd.PersistentFlags().StringVarP(&o.requestsMemory, "resources-requests-memory", "", "512M", "set resources request memory")
	cmd.PersistentFlags().StringVarP(&o.requestsMemory, "resources-requests-ephemeral-storage", "", "", "set resources request ephemeral-storage")
	return o
}

func (o *deployResourceOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployResourceOptions) complete() error {
	return nil
}

func (o *deployResourceOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployResourceOptions {
	if !o.enabled {
		return o
	}

	deployment.Spec.Resources = &kubemqcluster.ResourceConfig{
		LimitsCpu:                o.limitsCpu,
		LimitsMemory:             o.limitsMemory,
		LimitsEphemeralStorage:   o.limitsEphemeralStorage,
		RequestsCpu:              o.requestsCpu,
		RequestsMemory:           o.requestsMemory,
		RequestsEphemeralStorage: o.requestsEphemeralStorage,
	}
	return o
}
