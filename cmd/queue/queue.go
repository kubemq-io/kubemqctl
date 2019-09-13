package queue

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var queueExamples = `
	# Execute send to queue command
	kubemqctl queue send

	# Execute attached to a queue command
	kubemqctl queue attach

	# Execute receive to queue command
	kubemqctl queue receive
	
	# Execute list queue command
	kubemqctl queries list

	# Execute peek queue command
	kubemqctl queries peak

	# Execute ack queue command
	 kubemqctl queries ack

	# Execute stream queue command
	kubemqctl queries stream
`
var queueLong = `Execute KubeMQ 'queue' commands`
var queueShort = `Execute KubeMQ 'queue' commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdQueue(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "queue",
		Aliases: []string{"q", "qu"},
		Short:   queueShort,
		Long:    queueShort,
		Example: queueExamples,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdQueueSend(cfg))
	cmd.AddCommand(NewCmdQueueReceive(cfg))
	cmd.AddCommand(NewCmdQueuePeek(cfg))
	cmd.AddCommand(NewCmdQueueAck(cfg))
	cmd.AddCommand(NewCmdQueueList(cfg))
	cmd.AddCommand(NewCmdQueueStream(cfg))
	cmd.AddCommand(NewCmdQueueAttach(cfg))

	return cmd
}
