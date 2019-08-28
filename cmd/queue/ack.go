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

type QueueAckOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	wait      int
}

var queueAckExamples = `
	# ack all messages in a queue channel 'some-channel' with 2 seconds of wait to complete operation
	kubetools queue ack some-channel
	
	# ack all messages in a queue channel 'some-long-queue' with 30 seconds of wait to complete operation
	kubetools queue ack some-long-queue -w 30
`
var queueAckLong = `ack all messages in a queue`
var queueAckShort = `ack all messages in a queue`

func NewCmdQueueAck(cfg *config.Config, opts *QueueOptions) *cobra.Command {
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
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
		return fmt.Errorf("create send client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	res, err := client.AQM().
		SetChannel(o.channel).
		SetWaitTimeSeconds(o.wait).
		Send(ctx)
	if err != nil {
		return fmt.Errorf("ack all queue messages, %s", err.Error())
	}
	if res.IsError {
		return fmt.Errorf("peak queue message, %s", res.Error)
	}
	utils.Printlnf("acked %d messages", res.AffectedMessages)

	return nil
}
