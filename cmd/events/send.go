package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type EventsSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	message   string
	metadata  string
	iter      int
}

var eventsSendExamples = `
	# Send message to a events channel
	kubetools events send some-channel some-message
	
	# Send message to a events channel with metadata
	kubetools events send some-channel some-message -m some-metadata
	
	# Send 10 messages to a events channel
	kubetools events send some-channel some-message -i 10
`
var eventsSendLong = `send messages to a events channel`
var eventsSendShort = `send messages to a events channel`

func NewCmdEventsSend(cfg *config.Config, opts *EventsOptions) *cobra.Command {
	o := &EventsSendOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "send",
		Aliases: []string{"s"},
		Short:   eventsSendShort,
		Long:    eventsSendLong,
		Example: eventsSendExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "m", "", "set metadata message")
	cmd.PersistentFlags().IntVarP(&o.iter, "iterations", "i", 1, "set how many messages to send")

	return cmd
}

func (o *EventsSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 2 {
		o.channel = args[0]
		o.message = args[1]
		return nil
	}
	return fmt.Errorf("missing arguments, must be 2 arguments, channel and message")
}

func (o *EventsSendOptions) Validate() error {
	return nil
}

func (o *EventsSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()
	for i := 1; i <= o.iter; i++ {
		msg := client.E().
			SetChannel(o.channel).
			SetId(uuid.New().String()).
			SetBody([]byte(o.message)).
			SetMetadata(o.metadata)
		err = msg.Send(ctx)
		if err != nil {
			return fmt.Errorf("sending events message, %s", err.Error())
		}
		utils.Printlnf("[iteration: %d] [channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s}", i, msg.Channel, msg.ClientId, msg.Id, msg.Metadata, msg.Body)
	}
	return nil
}
