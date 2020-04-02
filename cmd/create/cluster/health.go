package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

type deployHealthOptions struct {
	enabled             bool
	initialDelaySeconds int32
	periodSeconds       int32
	timeoutSeconds      int32
	successThreshold    int32
	failureThreshold    int32
}

func setHealthOptions(cmd *cobra.Command) *deployHealthOptions {
	o := &deployHealthOptions{
		enabled: false,
	}
	cmd.PersistentFlags().BoolVarP(&o.enabled, "health-enabled", "", false, "enable resources configuration")
	cmd.PersistentFlags().Int32VarP(&o.initialDelaySeconds, "health-initial-delay", "", 10, "set health prob initial delay seconds ")
	cmd.PersistentFlags().Int32VarP(&o.periodSeconds, "health-period-seconds", "", 10, "set health prob period seconds ")
	cmd.PersistentFlags().Int32VarP(&o.timeoutSeconds, "health-timout-seconds", "", 5, "set health prob timeout seconds ")
	cmd.PersistentFlags().Int32VarP(&o.successThreshold, "health-success-threshold", "", 1, "set health prob success threshold")
	cmd.PersistentFlags().Int32VarP(&o.failureThreshold, "health-failure-threshold", "", 12, "set health prob failure threshold")

	return o
}

func (o *deployHealthOptions) validate() error {
	if !o.enabled {
		return nil
	}
	return nil
}

func (o *deployHealthOptions) complete() error {
	return nil
}

func (o *deployHealthOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployHealthOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Health = &kubemqcluster.HealthConfig{
		Enabled:             true,
		InitialDelaySeconds: o.initialDelaySeconds,
		PeriodSeconds:       o.periodSeconds,
		TimeoutSeconds:      o.timeoutSeconds,
		SuccessThreshold:    o.successThreshold,
		FailureThreshold:    o.failureThreshold,
	}

	return o
}
