package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"strings"
)

var serviceTypes = map[string]string{"clusterip": "ClusterIP", "nodeport": "NodePort", "loadbalancer": "LoadBalancer"}

var defaultApiConfig = &deployApiOptions{
	disabled: false,
	port:     8080,
	expose:   "ClusterIP",
	nodePort: 0,
}

type deployApiOptions struct {
	disabled bool
	port     int32
	expose   string
	nodePort int32
}

func setApiConfig(cmd *cobra.Command) *deployApiOptions {
	o := &deployApiOptions{
		disabled: false,
		port:     0,
		expose:   "",
		nodePort: 0,
	}
	cmd.PersistentFlags().BoolVarP(&o.disabled, "api-disabled", "", false, "disable Api interface")
	cmd.PersistentFlags().Int32VarP(&o.port, "api-port", "", 8080, "set api port value")
	cmd.PersistentFlags().StringVarP(&o.expose, "api-expose", "", "ClusterIP", "set api port service type (ClusterIP,NodePort,LoadBalancer)")
	cmd.PersistentFlags().Int32VarP(&o.nodePort, "api-node-port", "", 0, "set api node port value")
	return o
}

func (o *deployApiOptions) validate() error {
	if o.expose != "" {
		if _, ok := serviceTypes[strings.ToLower(o.expose)]; !ok {
			return fmt.Errorf("invalid api-expose value: %s", o.expose)
		}
	}

	return nil
}
func (o *deployApiOptions) complete() error {
	return nil
}

func (o *deployApiOptions) setConfig(deployment *cluster.KubemqCluster) *deployApiOptions {
	if isDefault(o, defaultApiConfig) {
		return o
	}
	deployment.Spec.Api = &cluster.ApiConfig{
		Disabled: o.disabled,
		Port:     o.port,
		Expose:   o.expose,
		NodePort: o.nodePort,
	}
	return o
}
