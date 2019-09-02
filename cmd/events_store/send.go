package events_store

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

type EventsStoreSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	message   string
	metadata  string
	iter      int
}

var eventsSendExamples = `
	# Send message to an events store channel
	kubetools events_store send some-channel some-message
	
	# Send message to an events store channel with metadata
	kubetools events_store send some-channel some-message -m some-metadata

	# Send 10 messages to an events store channel
	kubetools events_store send some-channel some-message -i 10
`
var eventsSendLong = `Send messages to an events store channel`
var eventsSendShort = `Send messages to an events store channel`

func NewCmdEventsStoreSend(cfg *config.Config) *cobra.Command {
	o := &EventsStoreSendOptions{
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
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "m", "", "Set metadata message")
	cmd.PersistentFlags().IntVarP(&o.iter, "iterations", "i", 1, "Set how many messages to send")

	return cmd
}

func (o *EventsStoreSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 2 {
		o.channel = args[0]
		o.message = args[1]
		return nil
	}
	return fmt.Errorf("missing arguments, must be 2 arguments, channel and a message")
}

func (o *EventsStoreSendOptions) Validate() error {
	return nil
}

func (o *EventsStoreSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()
	for i := 1; i <= o.iter; i++ {
		msg := client.ES().
			SetChannel(o.channel).
			SetId(uuid.New().String()).
			SetBody([]byte(o.message)).
			SetMetadata(o.metadata)
		res, err := msg.Send(ctx)
		if err != nil {
			return fmt.Errorf("sending events store message, %s", err.Error())
		}
		utils.Printlnf("[iteration: %d] [channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s, sent:%t}", i, msg.Channel, msg.ClientId, msg.Id, msg.Metadata, msg.Body, res.Sent)

	}
	return nil
}
