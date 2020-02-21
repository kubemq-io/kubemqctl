package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
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
	cmd.PersistentFlags().Int32VarP(&o.initialDelaySeconds, "health-initial-delay", "", 5, "set health prob initial delay seconds ")
	cmd.PersistentFlags().Int32VarP(&o.periodSeconds, "health-period", "", 10, "set health prob period seconds ")
	cmd.PersistentFlags().Int32VarP(&o.timeoutSeconds, "health-timout", "", 5, "set health prob timeout seconds ")
	cmd.PersistentFlags().Int32VarP(&o.successThreshold, "health-success", "", 1, "set health prob success threshold")
	cmd.PersistentFlags().Int32VarP(&o.failureThreshold, "health-failure", "", 6, "set health prob failure threshold")

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

func (o *deployHealthOptions) setConfig(deployment *cluster.KubemqCluster) *deployHealthOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Health = &cluster.HealthConfig{
		Enabled:             true,
		InitialDelaySeconds: o.initialDelaySeconds,
		PeriodSeconds:       o.periodSeconds,
		TimeoutSeconds:      o.timeoutSeconds,
		SuccessThreshold:    o.successThreshold,
		FailureThreshold:    o.failureThreshold,
	}

	return o
}
