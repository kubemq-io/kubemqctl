package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"time"
)

type QueueSendOptions struct {
	cfg        *config.Config
	transport  string
	expiration int
	delay      int
	channel    string
	message    string
	maxReceive int
	deadLetter string
}

var queueSendExamples = `
	# Send message to a queue
	kubetools queue send some-channel some-message

	# Send message to a queue with a message expiration of 5 seconds
	kubetools queue send some-channel some-message -e 5

	# Send message to a queue with a message delay of 5 seconds
	kubetools queue send some-channel some-message -d 5

	# Send message to a queue with a message policy of max receive 5 times and dead-letter queue 'dead-letter'
	kubetools queue send some-channel some-message -m 5 -q dead-letter
`
var queueSendLong = `send a message to a queue`
var queueSendShort = `send a message to a queue`

func NewCmdQueueSend(cfg *config.Config, opts *QueueOptions) *cobra.Command {
	o := &QueueSendOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "send",
		Aliases: []string{"s"},
		Short:   queueSendShort,
		Long:    queueSendLong,
		Example: queueSendExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().IntVarP(&o.expiration, "expiration", "e", 0, "set queue message expiration seconds")
	cmd.PersistentFlags().IntVarP(&o.delay, "delay", "d", 0, "set queue message send delay seconds")
	cmd.PersistentFlags().IntVarP(&o.maxReceive, "max-receive", "m", 0, "set dead-letter max receive count")
	cmd.PersistentFlags().StringVarP(&o.deadLetter, "dead-letter-queue", "q", "", "set dead-letter queue name")

	return cmd
}

func (o *QueueSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 2 {
		o.channel = args[0]
		o.message = args[1]
		return nil
	}
	return fmt.Errorf("missing arguments, must be 2 arguments, channel and message")
}

func (o *QueueSendOptions) Validate() error {
	return nil
}

func (o *QueueSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create send client, %s", err.Error())

	}
	defer utils.CheckErr(client.Close())
	res, err := client.QM().
		SetChannel(o.channel).
		SetBody([]byte(o.message)).
		SetPolicyExpirationSeconds(o.expiration).
		SetPolicyDelaySeconds(o.delay).
		SetPolicyMaxReceiveCount(o.maxReceive).
		SetPolicyMaxReceiveQueue(o.deadLetter).
		Send(ctx)
	if err != nil {
		return fmt.Errorf("sending queue message, %s", err.Error())
	}

	if res != nil {
		if res.IsError {
			return fmt.Errorf("sending queue message response, %s", res.Error)

		}
		utils.PrintfAndExit("queue message sent at: %s", time.Unix(0, res.SentAt))
	}
	return nil
}
