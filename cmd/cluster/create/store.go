package create

import (
	"github.com/kubemq-io/kubemqctl/pkg/k8s/crd/cluster"
	"github.com/spf13/cobra"
)

type deployStoreOptions struct {
	clean                    bool
	path                     string
	maxChannels              int32
	maxSubscribers           int32
	maxMessages              int32
	maxChannelSize           int32
	messagesRetentionMinutes int32
	purgeInactiveMinutes     int32
}

func defaultStoreConfig(cmd *cobra.Command) *deployStoreOptions {
	o := &deployStoreOptions{
		clean:                    false,
		path:                     "",
		maxChannels:              0,
		maxSubscribers:           0,
		maxMessages:              0,
		maxChannelSize:           0,
		messagesRetentionMinutes: 0,
		purgeInactiveMinutes:     0,
	}

	cmd.PersistentFlags().BoolVarP(&o.clean, "store-clean", "", false, "set clear persistence data on start-up    ")
	cmd.PersistentFlags().StringVarP(&o.path, "store-path", "", "./store", "set persistence file path")
	cmd.PersistentFlags().Int32VarP(&o.maxChannels, "store-max-channels", "", 0, "set limit number of persistence channels")
	cmd.PersistentFlags().Int32VarP(&o.maxSubscribers, "store-max-subscribers", "", 0, "set limit of subscribers per channel")
	cmd.PersistentFlags().Int32VarP(&o.maxMessages, "store-max-messages", "", 0, "set limit of messages per channel")
	cmd.PersistentFlags().Int32VarP(&o.maxChannelSize, "store-max-channel-size", "", 0, "Set limit size of channel in bytes ")
	cmd.PersistentFlags().Int32VarP(&o.messagesRetentionMinutes, "store-messages-retention-minutes", "", 1440, "set message retention time in minutes")
	cmd.PersistentFlags().Int32VarP(&o.purgeInactiveMinutes, "store-purge-inactive-minutes", "", 1440, "set time in minutes of channel inactivity to delete")

	return o
}

func (o *deployStoreOptions) validate() error {
	return nil
}
func (o *deployStoreOptions) complete() error {
	return nil
}

func (o *deployStoreOptions) setConfig(deployment *cluster.KubemqCluster) *deployStoreOptions {
	deployment.Spec.Store = &cluster.StoreConfig{
		Clean:                    o.clean,
		Path:                     o.path,
		MaxChannels:              new(int32),
		MaxSubscribers:           new(int32),
		MaxMessages:              new(int32),
		MaxChannelSize:           new(int32),
		MessagesRetentionMinutes: new(int32),
		PurgeInactiveMinutes:     new(int32),
	}

	*deployment.Spec.Store.MaxChannels = o.maxChannels
	*deployment.Spec.Store.MaxSubscribers = o.maxSubscribers
	*deployment.Spec.Store.MaxMessages = o.maxMessages
	*deployment.Spec.Store.MaxChannelSize = o.maxChannelSize
	*deployment.Spec.Store.MessagesRetentionMinutes = o.messagesRetentionMinutes
	*deployment.Spec.Store.PurgeInactiveMinutes = o.purgeInactiveMinutes
	return o
}
