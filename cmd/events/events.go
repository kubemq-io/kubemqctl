package events

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

type EventsOptions struct {
	transport string
}

var eventsExamples = ``
var eventsLong = ``
var eventsShort = `Execute KubeMQ events commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdEvents(cfg *config.Config) *cobra.Command {
	o := &EventsOptions{
		transport: "grpc",
	}
	cmd := &cobra.Command{
		Use:     "events",
		Aliases: []string{"e", ""},
		Short:   eventsShort,
		Long:    eventsShort,
		Example: eventsExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdEventsSend(cfg, o))
	cmd.AddCommand(NewCmdEventsReceive(cfg, o))
	cmd.AddCommand(NewCmdEventsAttach(cfg, o))

	cmd.PersistentFlags().StringVarP(&o.transport, "transport", "t", "grpc", "set transport type, grpc or rest")
	return cmd
}
