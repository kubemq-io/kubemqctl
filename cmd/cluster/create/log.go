package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

type deployLogOptions struct {
	level int32
	file  string
}

func defaultLogConfig(cmd *cobra.Command) *deployLogOptions {
	o := &deployLogOptions{
		level: 2,
		file:  "",
	}
	cmd.PersistentFlags().Int32VarP(&o.level, "log-data", "", 2, "set log level")
	cmd.PersistentFlags().StringVarP(&o.file, "log-file", "", "", "set log filename")
	return o
}

func (o *deployLogOptions) validate() error {
	return nil
}
func (o *deployLogOptions) complete() error {
	return nil
}

func (o *deployLogOptions) setConfig(deployment *cluster.KubemqCluster) *deployLogOptions {
	deployment.Spec.Log = &cluster.LogConfig{
		Level: new(int32),
		File:  o.file,
	}
	*deployment.Spec.Log.Level = o.level
	return o
}
