package events_store

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var eventsExamples = `
	# Execute send events_store command 
	# kubemqctl events_store send

	# Execute receive events_store command
	# kubemqctl events_store receive

	# Execute attach to events_store command
	# kubemqctl events_store attach

	# Execute list of events_store channels command
`
var eventsLong = `Execute KubeMQ 'events_store' Pub/Sub commands`
var eventsShort = `Execute KubeMQ 'events_store' Pub/Sub commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdEventsStore(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "events_store",
		Aliases:   []string{"es"},
		Short:     eventsLong,
		Long:      eventsShort,
		Example:   eventsExamples,
		ValidArgs: []string{"send", "receive", "attach", "list"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdEventsStoreSend(cfg))
	cmd.AddCommand(NewCmdEventsStoreReceive(cfg))
	cmd.AddCommand(NewCmdEventsStoreAttach(cfg))
	cmd.AddCommand(NewCmdEventsStoreList(cfg))

	return cmd
}
