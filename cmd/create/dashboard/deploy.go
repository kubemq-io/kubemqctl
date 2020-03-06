package dashboard

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqdashboard"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployOptions struct {
	name       string
	namespace  string
	port       int32
	prometheus *deployPrometheusOptions
	grafana    *deployGrafanaOptions
}

func defaultDeployOptions(cmd *cobra.Command) *deployOptions {
	o := &deployOptions{
		port:       0,
		prometheus: setPrometheusConfig(cmd),
		grafana:    setGrafanaConfig(cmd),
	}
	cmd.PersistentFlags().StringVarP(&o.name, "name", "", "kubemq-dashboard", "set kubemq dashboard name")
	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "n", "kubemq", "set kubemq dashboard namespace")
	cmd.PersistentFlags().Int32VarP(&o.port, "port", "p", 32000, "set kubemq dashboard port")
	return o
}

func (o *deployOptions) validate() error {

	if err := o.prometheus.validate(); err != nil {
		return err
	}

	if err := o.grafana.validate(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) complete() error {

	if err := o.prometheus.complete(); err != nil {
		return err
	}

	if err := o.grafana.complete(); err != nil {
		return err
	}
	return nil
}

func (o *deployOptions) getDashboardDeployment() *kubemqdashboard.KubemqDashboard {

	deployment := &kubemqdashboard.KubemqDashboard{
		TypeMeta: v1.TypeMeta{
			Kind:       "KubemqDashboard",
			APIVersion: "core.k8s.kubemq.io/v1alpha1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      o.name,
			Namespace: o.namespace,
		},
		Spec: kubemqdashboard.KubemqDashboardSpec{
			Port:       o.port,
			Prometheus: nil,
			Grafana:    nil,
		},
		Status: kubemqdashboard.KubemqDashboardStatus{},
	}
	o.prometheus.setConfig(deployment)
	o.grafana.setConfig(deployment)

	return deployment
}
