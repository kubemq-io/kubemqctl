package dashboard

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
	"github.com/spf13/cobra"
)

var defaultPrometheusOptions = &deployPrometheusOptions{
	nodePort: 0,
	image:    "",
}

type deployPrometheusOptions struct {
	nodePort int32
	image    string
}

func setPrometheusConfig(cmd *cobra.Command) *deployPrometheusOptions {
	o := &deployPrometheusOptions{
		nodePort: 0,
		image:    "",
	}
	cmd.PersistentFlags().Int32VarP(&o.nodePort, "prometheus-port", "", 0, "set export prometheus port")
	cmd.PersistentFlags().StringVarP(&o.image, "prometheus-image", "", "", "set prometheus docker image")
	return o
}

func (o *deployPrometheusOptions) validate() error {
	return nil
}
func (o *deployPrometheusOptions) complete() error {
	return nil
}

func (o *deployPrometheusOptions) setConfig(deployment *kubemqdashboard.KubemqDashboard) *deployPrometheusOptions {
	if isDefault(o, defaultPrometheusOptions) {
		return o
	}

	deployment.Spec.Prometheus = &kubemqdashboard.PrometheusConfig{
		NodePort: o.nodePort,
		Image:    o.image,
	}
	return o
}
