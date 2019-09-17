package queue

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var queueExamples = `
	# Execute send to queue command
	kubemqctl queues send

	# Execute attached to a queue command
	kubemqctl queues attach

	# Execute receive to queue command
	kubemqctl queues receive
	
	# Execute list queue command
	kubemqctl queues list

	# Execute peek queue command
	kubemqctl queues peak

	# Execute ack queue command
	 kubemqctl queues ack

	# Execute stream queue command
	kubemqctl queues stream
`
var queueLong = `Execute KubeMQ 'queues' commands`
var queueShort = `Execute KubeMQ 'queues' commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdQueue(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:       "queues",
		Aliases:   []string{"q", "qu", "queue"},
		Short:     queueShort,
		Long:      queueShort,
		Example:   queueExamples,
		ValidArgs: []string{"send", "receive", "attach", "peek", "ack", "list", "stream"},
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
