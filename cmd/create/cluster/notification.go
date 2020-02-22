package cluster

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/types/kubemqcluster"
	"github.com/spf13/cobra"
)

type deployNotificationOptions struct {
	enabled bool
	prefix  string
	log     bool
}

func setNotificationConfig(cmd *cobra.Command) *deployNotificationOptions {
	o := &deployNotificationOptions{
		enabled: false,
		prefix:  "",
		log:     false,
	}

	cmd.PersistentFlags().BoolVarP(&o.enabled, "notification-enabled", "", false, "set notification enable")
	cmd.PersistentFlags().StringVarP(&o.prefix, "notification-prefix", "", "", "set notification channel prefix")
	cmd.PersistentFlags().BoolVarP(&o.log, "notification-log", "", false, "set log notification to std-out")

	return o
}

func (o *deployNotificationOptions) validate() error {

	return nil
}
func (o *deployNotificationOptions) complete() error {
	return nil
}

func (o *deployNotificationOptions) setConfig(deployment *kubemqcluster.KubemqCluster) *deployNotificationOptions {
	if !o.enabled {
		return o
	}
	deployment.Spec.Notification = &kubemqcluster.NotificationConfig{
		Enabled: o.enabled,
		Prefix:  o.prefix,
		Log:     o.log,
	}
	return o
}
