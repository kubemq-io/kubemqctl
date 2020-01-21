package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
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

func defaultHealthOptions(cmd *cobra.Command) *deployHealthOptions {
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

func (o *deployHealthOptions) setConfig(config *deployment.KubeMQManifestConfig) *deployHealthOptions {
	if !o.enabled {
		return o
	}
	tmpl := `          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: %d
            periodSeconds: %d
            timeoutSeconds: %d
            successThreshold: %d
            failureThreshold: %d
`
	prob := fmt.Sprintf(tmpl,
		o.initialDelaySeconds,
		o.timeoutSeconds,
		o.periodSeconds,
		o.successThreshold,
		o.failureThreshold)
	config.StatefulSet.SetHealthProbe(prob)
	return o
}
