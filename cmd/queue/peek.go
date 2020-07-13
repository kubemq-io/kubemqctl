package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type QueuePeekOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	messages  int
	wait      int
}

var queuePeekExamples = `
	# Peek 1 messages from a queue and wait for 2 seconds (default)
	kubemqctl queue peek some-channel

	# Peek 3 messages from a queue and wait for 5 seconds
	kubemqctl queue peek some-channel -m 3 -w 5
`
var queuePeekLong = `Peek command allows to peek one or many messages from a queue channel without removing them from the queue`
var queuePeekShort = `Peek a messages from a queue channel command`

func NewCmdQueuePeek(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueuePeekOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "peek",
		Aliases: []string{"p"},
		Short:   queuePeekShort,
		Long:    queuePeekLong,
		Example: queuePeekExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set how many messages we want to peek from queue")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait", "w", 2, "set how many seconds to wait for peeking queue messages")

	return cmd
}

func (o *QueuePeekOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueuePeekOptions) Validate() error {
	return nil
}

func (o *QueuePeekOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	res, err := client.RQM().
		SetChannel(o.channel).
		SetWaitTimeSeconds(o.wait).
		SetMaxNumberOfMessages(o.messages).
		SetIsPeak(true).
		Send(ctx)
	if err != nil {
		return fmt.Errorf("peek queue message, %s", err.Error())
	}
	if res.IsError {
		return fmt.Errorf("peek queue message, %s", res.Error)
	}

	if res.MessagesReceived > 0 {
		utils.Printlnf("peeking %d messages", res.MessagesReceived)
		printItems(res.Messages)
	} else {
		utils.Printlnf("no messages in queue to peek")
	}

	return nil
}
