package events

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var eventsExamples = `
	# Execute send events command
	# kubemqctl events send

	# Execute receive an events command
	# kubemqctl events receive

	# Execute attach to an events command
	# kubemqctl events attach

`
var eventsLong = `Execute KubeMQ 'events' Pub/Sub commands`
var eventsShort = `Execute KubeMQ 'events' Pub/Sub commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdEvents(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "events",
		Aliases:   []string{"e"},
		Short:     eventsShort,
		Long:      eventsLong,
		Example:   eventsExamples,
		ValidArgs: []string{"send", "receive", "attach"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdEventsSend(cfg))
	cmd.AddCommand(NewCmdEventsReceive(cfg))
	cmd.AddCommand(NewCmdEventsAttach(cfg))

	return cmd
}
