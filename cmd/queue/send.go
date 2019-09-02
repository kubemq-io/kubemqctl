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
	metadata   string
	deadLetter string
}

var queueSendExamples = `
	# Send message to a queue channel
	kubetools queue send some-channel some-message
	
	# Send message to a queue channel with metadata
	kubetools queue send some-channel some-message -m some-metadata
	
	# Send message to a queue with a message expiration of 5 seconds
	kubetools queue send some-channel some-message -e 5

	# Send message to a queue with a message delay of 5 seconds
	kubetools queue send some-channel some-message -d 5

	# Send message to a queue with a message policy of max receive 5 times and dead-letter queue 'dead-letter'
	kubetools queue send some-channel some-message -r 5 -q dead-letter
`
var queueSendLong = `Send a message to a queue channel`
var queueSendShort = `Send a message to a queue channel`

func NewCmdQueueSend(cfg *config.Config) *cobra.Command {
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
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().IntVarP(&o.expiration, "expiration", "e", 0, "Set queue message expiration seconds")
	cmd.PersistentFlags().IntVarP(&o.delay, "delay", "d", 0, "Set queue message send delay seconds")
	cmd.PersistentFlags().IntVarP(&o.maxReceive, "max-receive", "r", 0, "Set dead-letter max receive count")
	cmd.PersistentFlags().StringVarP(&o.deadLetter, "dead-letter-queue", "q", "", "Set dead-letter queue name")
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "m", "", "Set metadata message")

	return cmd
}

func (o *QueueSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 2 {
		o.channel = args[0]
		o.message = args[1]
		return nil
	}
	return fmt.Errorf("missing arguments, must be 2 arguments, channel and a message")
}

func (o *QueueSendOptions) Validate() error {
	return nil
}

func (o *QueueSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}

	defer func() {
		client.Close()
	}()
	msg := client.QM().
		SetChannel(o.channel).
		SetBody([]byte(o.message)).
		SetMetadata(o.metadata).
		SetPolicyExpirationSeconds(o.expiration).
		SetPolicyDelaySeconds(o.delay).
		SetPolicyMaxReceiveCount(o.maxReceive).
		SetPolicyMaxReceiveQueue(o.deadLetter)
	res, err := msg.Send(ctx)
	if err != nil {
		return fmt.Errorf("sending queue message, %s", err.Error())
	}

	if res != nil {
		if res.IsError {
			return fmt.Errorf("sending queue message response, %s", res.Error)

		}
		var delay string
		var exp string
		if res.DelayedTo > 0 {
			delay = fmt.Sprintf(", delayed to: %s", time.Unix(0, res.DelayedTo))
		}
		if res.ExpirationAt > 0 {
			exp = fmt.Sprintf(", expire at: %s", time.Unix(0, res.ExpirationAt))
		}
		utils.Printlnf("[channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s, sent at: %s%s%s}", msg.Channel, msg.ClientId, res.MessageID, msg.Metadata, msg.Body, time.Unix(0, res.SentAt).Format("2006-01-02 15:04:05"), exp, delay)
	}
	return nil
}
