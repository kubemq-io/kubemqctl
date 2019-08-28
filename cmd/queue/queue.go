package queue

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

type QueueOptions struct {
	transport string
}

var queueExamples = ``
var queueLong = ``
var queueShort = `Execute KubeMQ queue commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdQueue(cfg *config.Config) *cobra.Command {
	o := &QueueOptions{
		transport: "grpc",
	}
	cmd := &cobra.Command{
		Use:     "queue",
		Aliases: []string{"q", "qu"},
		Short:   queueShort,
		Long:    queueShort,
		Example: queueExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdQueueSend(cfg, o))
	cmd.AddCommand(NewCmdQueueReceive(cfg, o))
	cmd.AddCommand(NewCmdQueuePeak(cfg, o))
	cmd.AddCommand(NewCmdQueueAck(cfg, o))
	cmd.AddCommand(NewCmdQueueList(cfg, o))
	cmd.AddCommand(NewCmdQueueStream(cfg, o))
	cmd.AddCommand(NewCmdQueueAttach(cfg, o))

	cmd.PersistentFlags().StringVarP(&o.transport, "transport", "t", "grpc", "set transport type, grpc or rest")
	return cmd
}
