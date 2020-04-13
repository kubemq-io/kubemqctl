package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

var defaultLogConfig = &deployLogOptions{
	level: 2,
	file:  "",
}

type deployLogOptions struct {
	level int32
	file  string
}

func setLogConfig(cmd *cobra.Command) *deployLogOptions {
	o := &deployLogOptions{
		level: 2,
		file:  "",
	}
	cmd.PersistentFlags().Int32VarP(&o.level, "log-level", "", 2, "set log level")
	cmd.PersistentFlags().StringVarP(&o.file, "log-file", "", "", "set log filename")
	return o
}

func (o *deployLogOptions) validate() error {
	return nil
}
func (o *deployLogOptions) complete() error {
	return nil
}

func (o *deployLogOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployLogOptions {
	if isDefault(o, defaultLogConfig) {
		return o
	}
	deployment.Spec.Log = &kubemqcluster.LogConfig{
		Level: new(int32),
		File:  o.file,
	}
	*deployment.Spec.Log.Level = o.level
	return o
}
