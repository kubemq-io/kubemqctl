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

type QueueAckOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	wait      int
}

var queueAckExamples = `
	# Ack all messages in a 'queues' channel 'some-channel' with 2 seconds of wait to complete operation
	kubemqctl queue ack some-channel
	
	# Ack all messages in a 'queues' channel 'some-long-queue' with 30 seconds of wait to complete operation
	kubemqctl queue ack some-long-queue -w 30
`
var queueAckLong = `Ack command allows to clear / remove / ack all messages in a 'queues' channel`
var queueAckShort = `Ack all messages in a 'queues' channel`

func NewCmdQueueAck(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueueAckOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "ack",
		Aliases: []string{"ac"},
		Short:   queueAckShort,
		Long:    queueAckLong,
		Example: queueAckExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.wait, "wait", "w", 2, "set how many seconds to wait for ack all messages")

	return cmd
}

func (o *QueueAckOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueueAckOptions) Validate() error {
	return nil
}

func (o *QueueAckOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	res, err := client.AQM().
		SetChannel(o.channel).
		SetWaitTimeSeconds(o.wait).
		Send(ctx)
	if err != nil {
		return fmt.Errorf("ack all 'queues' messages, %s", err.Error())
	}
	if res.IsError {
		return fmt.Errorf("ack all 'queues' message, %s", res.Error)
	}
	utils.Printlnf("acked %d messages", res.AffectedMessages)

	return nil
}
