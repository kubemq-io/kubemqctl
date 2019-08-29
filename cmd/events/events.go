package events

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

var eventsExamples = ``
var eventsLong = ``
var eventsShort = `Execute KubeMQ events commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdEvents(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "events",
		Aliases: []string{"e"},
		Short:   eventsShort,
		Long:    eventsShort,
		Example: eventsExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdEventsSend(cfg))
	cmd.AddCommand(NewCmdEventsReceive(cfg))
	cmd.AddCommand(NewCmdEventsAttach(cfg))
	cmd.AddCommand(NewCmdEventsMetrics(cfg))

	return cmd
}
