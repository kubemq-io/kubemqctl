package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

type deployVolumeOptions struct {
	size string
}

func defaultVolumeConfig(cmd *cobra.Command) *deployVolumeOptions {
	o := &deployVolumeOptions{
		size: "",
	}
	cmd.PersistentFlags().StringVarP(&o.size, "volume-size", "", "", "set persisted volume size")
	return o
}

func (o *deployVolumeOptions) validate() error {
	return nil
}
func (o *deployVolumeOptions) complete() error {
	return nil
}

func (o *deployVolumeOptions) setConfig(deployment *cluster.KubemqCluster) *deployVolumeOptions {
	deployment.Spec.Volume = &cluster.VolumeConfig{
		Size: o.size,
	}
	return o
}
