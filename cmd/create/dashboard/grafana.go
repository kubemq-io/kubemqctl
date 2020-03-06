package dashboard

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
	"github.com/spf13/cobra"
)

var defaultGrafanaOptions = &deployGrafanaOptions{
	dashboardUrl: "",
	image:        "",
}

type deployGrafanaOptions struct {
	dashboardUrl string
	image        string
}

func setGrafanaConfig(cmd *cobra.Command) *deployGrafanaOptions {
	o := &deployGrafanaOptions{
		dashboardUrl: "",
		image:        "",
	}

	cmd.PersistentFlags().StringVarP(&o.image, "grafana-image", "", "", "set grafana docker image")
	cmd.PersistentFlags().StringVarP(&o.dashboardUrl, "grafana-dashboard-url", "", "", "set grafana dashboard url image")
	return o
}

func (o *deployGrafanaOptions) validate() error {
	return nil
}
func (o *deployGrafanaOptions) complete() error {
	return nil
}

func (o *deployGrafanaOptions) setConfig(deployment *kubemqdashboard.KubemqDashboard) *deployGrafanaOptions {
	if isDefault(o, defaultGrafanaOptions) {
		return o
	}

	deployment.Spec.Grafana = &kubemqdashboard.GrafanaConfig{
		DashboardUrl: o.dashboardUrl,
		Image:        o.image,
	}
	return o
}
