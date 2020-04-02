package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

var defaultVolumeConfig = &deployVolumeOptions{
	size: "",
}

type deployVolumeOptions struct {
	size         string
	storageClass string
}

func setVolumeConfig(cmd *cobra.Command) *deployVolumeOptions {
	o := &deployVolumeOptions{
		size:         "",
		storageClass: "",
	}
	cmd.PersistentFlags().StringVarP(&o.size, "volume-size", "v", "", "set persisted volume size")
	cmd.PersistentFlags().StringVarP(&o.storageClass, "volume-storage-class", "", "", "set persisted volume storage class")
	return o
}

func (o *deployVolumeOptions) validate() error {
	return nil
}
func (o *deployVolumeOptions) complete() error {
	return nil
}

func (o *deployVolumeOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployVolumeOptions {
	if isDefault(o, defaultVolumeConfig) {
		return o
	}
	deployment.Spec.Volume = &kubemqcluster.VolumeConfig{
		Size:         o.size,
		StorageClass: o.storageClass,
	}
	return o
}
