package create

import (
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
	"strings"
)

type deployGrpcOptions struct {
	disabled   bool
	port       int32
	expose     string
	nodePort   int32
	bufferSize int32
	BodyLimit  int32
}

func defaultGrpcConfig(cmd *cobra.Command) *deployGrpcOptions {
	o := &deployGrpcOptions{
		disabled: false,
		port:     0,
		expose:   "",
		nodePort: 0,
	}
	cmd.PersistentFlags().BoolVarP(&o.disabled, "grpc-disabled", "", false, "disable grpc interface")
	cmd.PersistentFlags().Int32VarP(&o.port, "grpc-port", "", 50000, "set grpc port value")
	cmd.PersistentFlags().StringVarP(&o.expose, "grpc-expose", "", "ClusterIP", "set grpc port service type (ClusterIP,NodePort,LoadBalancer)")
	cmd.PersistentFlags().Int32VarP(&o.nodePort, "grpc-node-port", "", 0, "set grpc node port value")
	cmd.PersistentFlags().Int32VarP(&o.bufferSize, "grpc-buffer-size", "", 0, "set subscribe message / requests buffer size to use on server ")
	cmd.PersistentFlags().Int32VarP(&o.BodyLimit, "grpc-body-limit", "", 0, "set Max size of payload in bytes ")
	return o
}
func (o *deployGrpcOptions) validate() error {
	if o.expose != "" {
		if _, ok := serviceTypes[strings.ToLower(o.expose)]; !ok {
			return fmt.Errorf("invalid grpc-expose value: %s", o.expose)
		}
	}

	return nil
}
func (o *deployGrpcOptions) complete() error {
	return nil
}

func (o *deployGrpcOptions) setConfig(deployment *cluster.KubemqCluster) *deployGrpcOptions {
	deployment.Spec.Grpc = &cluster.GrpcConfig{
		Disabled:   o.disabled,
		Port:       o.port,
		Expose:     o.expose,
		NodePort:   o.nodePort,
		BufferSize: o.bufferSize,
		BodyLimit:  o.BodyLimit,
	}
	return o
}
