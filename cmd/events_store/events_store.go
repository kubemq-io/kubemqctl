package events_store

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

type EventsStoreOptions struct {
	transport string
}

var eventsExamples = ``
var eventsLong = ``
var eventsShort = `Execute KubeMQ events store commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdEventsStore(cfg *config.Config) *cobra.Command {
	o := &EventsStoreOptions{
		transport: "grpc",
	}
	cmd := &cobra.Command{
		Use:     "events_store",
		Aliases: []string{"es"},
		Short:   eventsShort,
		Long:    eventsShort,
		Example: eventsExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdEventsStoreSend(cfg, o))
	cmd.AddCommand(NewCmdEventsStoreReceive(cfg, o))
	cmd.AddCommand(NewCmdEventsStoreAttach(cfg, o))
	cmd.AddCommand(NewCmdEventsStoreList(cfg, o))

	cmd.PersistentFlags().StringVarP(&o.transport, "transport", "t", "grpc", "set transport type, grpc or rest")
	return cmd
}
