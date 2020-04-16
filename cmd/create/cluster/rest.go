package cluster

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
	"strings"
)

var defaultRestConfig = &deployRestOptions{
	disabled:   false,
	port:       9090,
	expose:     "ClusterIP",
	nodePort:   0,
	bufferSize: 0,
	BodyLimit:  0,
}

type deployRestOptions struct {
	disabled   bool
	port       int32
	expose     string
	nodePort   int32
	bufferSize int32
	BodyLimit  int32
}

func setRestConfig(cmd *cobra.Command) *deployRestOptions {
	o := &deployRestOptions{
		disabled: false,
		port:     0,
		expose:   "",
		nodePort: 0,
	}
	cmd.PersistentFlags().BoolVarP(&o.disabled, "rest-disabled", "", false, "disable rest interface")
	cmd.PersistentFlags().Int32VarP(&o.port, "rest-port", "", 9090, "set rest port value")
	cmd.PersistentFlags().StringVarP(&o.expose, "rest-expose", "", "ClusterIP", "set rest port service type (ClusterIP,NodePort,LoadBalancer)")
	cmd.PersistentFlags().Int32VarP(&o.nodePort, "rest-node-port", "", 0, "set rest node port value")
	cmd.PersistentFlags().Int32VarP(&o.bufferSize, "rest-buffer-size", "", 0, "set subscribe message / requests buffer size to use on server ")
	cmd.PersistentFlags().Int32VarP(&o.BodyLimit, "rest-body-limit", "", 0, "set Max size of payload in bytes ")
	return o
}
func (o *deployRestOptions) validate() error {
	if o.expose != "" {
		if _, ok := serviceTypes[strings.ToLower(o.expose)]; !ok {
			return fmt.Errorf("invalid rest-expose value: %s", o.expose)
		}
	}

	return nil
}
func (o *deployRestOptions) complete() error {
	return nil
}

func (o *deployRestOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployRestOptions {
	if isDefault(o, defaultRestConfig) {
		return o
	}
	deployment.Spec.Rest = &kubemqcluster.RestConfig{
		Disabled:   o.disabled,
		Port:       o.port,
		Expose:     o.expose,
		NodePort:   o.nodePort,
		BufferSize: o.bufferSize,
		BodyLimit:  o.BodyLimit,
	}
	return o
}
