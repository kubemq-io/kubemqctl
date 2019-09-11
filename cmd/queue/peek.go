package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
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
	kubetools queue peek some-channel

	# Peek 3 messages from a queue and wait for 5 seconds
	kubetools queue peek some-channel -m 3 -w 5
`
var queuePeekLong = `Peek a messages from a queue channel`
var queuePeekShort = `Peek a messages from a queue channel`

func NewCmdQueuePeek(cfg *config.Config) *cobra.Command {
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "Set how many messages we want to get from queue")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait", "w", 2, "Set how many seconds to wait for queue messages")

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
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
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
