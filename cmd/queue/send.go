package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/targets"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type QueueSendOptions struct {
	cfg        *config.Config
	transport  string
	expiration int
	delay      int
	channel    string
	body       string
	maxReceive int
	metadata   string
	deadLetter string
	messages   int
	fileName   string
	build      bool
}

var queueSendExamples = `
	# Send message to a queue channel channel
	kubemqctl queue send q1 some-message
	
	# Send message to a queue channel with metadata
	kubemqctl queue send q1 some-message --metadata some-metadata
	
	# Send 5 messages to a queues channel with metadata
	kubemqctl queue send q1 some-message --metadata some-metadata -m 5
	
	# Send message to a queue channel with a message expiration of 5 seconds
	kubemqctl queue send q1 some-message -e 5

	# Send message to a queue channel with a message delay of 5 seconds
	kubemqctl queue send q1 some-message -d 5

	# Send message to a queue channel with a message policy of max receive 5 times and dead-letter queue 'dead-letter'
	kubemqctl queue send q1 some-message -r 5 -q dead-letter
`
var queueSendLong = `Send command allows to send one or many message to a queue channel`
var queueSendShort = `Send a message to a queue channel command`

func NewCmdQueueSend(ctx context.Context, cfg *config.Config) *cobra.Command {
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
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().IntVarP(&o.expiration, "expiration", "e", 0, "set queue message expiration seconds")
	cmd.PersistentFlags().IntVarP(&o.delay, "delay", "d", 0, "set queue message sending delay seconds")
	cmd.PersistentFlags().IntVarP(&o.maxReceive, "max-receive", "r", 0, "set dead-letter max receive count")
	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set dead-letter max receive count")
	cmd.PersistentFlags().StringVarP(&o.deadLetter, "dead-letter-queue", "q", "", "set dead-letter queue name")
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "", "", "set queue message metadata field")
	cmd.PersistentFlags().StringVarP(&o.fileName, "file", "f", "", "set load message body from file")
	cmd.PersistentFlags().BoolVarP(&o.build, "build", "b", false, "build kubemq targets request")

	return cmd
}

func (o *QueueSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]

	} else {
		return fmt.Errorf("missing channel argument")
	}
	if o.build {
		data, err := targets.BuildRequest()
		if err != nil {
			return err
		}
		o.body = string(data)
		return nil
	}
	if o.fileName != "" {
		data, err := ioutil.ReadFile(o.fileName)
		if err != nil {
			return err
		}
		o.body = string(data)
	} else {
		if len(args) >= 2 {
			o.body = args[1]
		} else {
			return fmt.Errorf("missing body argument")
		}
	}
	return nil

}

func (o *QueueSendOptions) Validate() error {
	return nil
}

func (o *QueueSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()

	for i := 0; i < o.messages; i++ {
		msg := client.QM().
			SetChannel(o.channel).
			SetBody([]byte(o.body)).
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
			fmt.Println("Sent:")
			printQueueMessage(msg)
		}
	}

	return nil
}
