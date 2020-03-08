package events

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	kubemq2 "github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"time"
)

type EventsSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	message   string
	metadata  string
	messages  int
	isStream  bool
}

var eventsSendExamples = `
	# Send (Publish) message to a 'events' channel
	kubemqctl events send some-channel some-message
	
	# Send (Publish) message to a 'events' channel with metadata
	kubemqctl events send some-channel some-message --metadata some-metadata
	
	# Send (Publish) batch of 10 messages to a 'events' channel
	kubemqctl events send some-channel some-message -m 10
`
var eventsSendLong = `Send command allows to send (publish) one or many messages to an 'events' channel`
var eventsSendShort = `Send messages to an 'events' channel command`

func NewCmdEventsSend(ctx context.Context, cfg *config.Config) *cobra.Command {
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
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "", "", "set message metadata field")
	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set how many 'events' messages to send")
	cmd.PersistentFlags().BoolVarP(&o.isStream, "stream", "s", false, "set stream of all messages at once")

	return cmd
}

func (o *EventsSendOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 2 {
		o.channel = args[0]
		o.message = args[1]
		return nil
	}
	return fmt.Errorf("missing arguments, must be 2 arguments, channel and a message")
}

func (o *EventsSendOptions) Validate() error {
	return nil
}

func (o *EventsSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()

	if o.isStream {
		utils.Printlnf("Streaming %d events messages ...", o.messages)
		eventsCh := make(chan *kubemq2.Event, 100)
		errCh := make(chan error, 10)

		go client.StreamEvents(ctx, eventsCh, errCh)
		startTime := time.Now()
		for i := 1; i <= o.messages; i++ {
			eventsCh <- client.E().
				SetChannel(o.channel).
				SetId(uuid.New().String()).
				SetBody([]byte(o.message)).
				SetMetadata(o.metadata)
		}
		utils.Printlnf("%d events messages streamed in %s.", o.messages, time.Since(startTime))
		time.Sleep(time.Second)
	} else {
		for i := 1; i <= o.messages; i++ {
			msg := client.E().
				SetChannel(o.channel).
				SetId(uuid.New().String()).
				SetBody([]byte(o.message)).
				SetMetadata(o.metadata)
			err = msg.Send(ctx)
			if err != nil {
				return fmt.Errorf("sending 'events' message, %s", err.Error())
			}
			utils.Printlnf("[message: %d] [channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s}", i, msg.Channel, msg.ClientId, msg.Id, msg.Metadata, msg.Body)
		}
	}

	return nil
}
