package events_store

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	kubemq2 "github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/targets"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"time"
)

type EventsStoreSendOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	body      string
	metadata  string
	messages  int
	isStream  bool
	fileName  string
	build     bool
}

var eventsSendExamples = `
	# Send (Publish) body to an 'events store' channel
	kubemqctl events_store send some-channel some-body
	
	# Send (Publish) body to an 'events store' channel with metadata
	kubemqctl events_store send some-channel some-body --metadata some-metadata

	# Send 10 messages to an 'events store' channel
	kubemqctl events_store send some-channel some-body -m 10

	# Send 100 messages to an 'events store' channel in stream mode
	kubemqctl events_store send some-channel some-body -m 100 -s
`
var eventsSendLong = `Send command allows to send (publish) one or many messages to an 'events store' channel`
var eventsSendShort = `Send messages to an 'events store' channel command`

func NewCmdEventsStoreSend(ctx context.Context, cfg *config.Config) *cobra.Command {
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
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.metadata, "metadata", "", "", "set body metadata field")
	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set how many 'events store' messages to send")
	cmd.PersistentFlags().BoolVarP(&o.isStream, "stream", "s", false, "set stream of all messages at once")
	cmd.PersistentFlags().StringVarP(&o.fileName, "file", "f", "", "set load body from file")
	cmd.PersistentFlags().BoolVarP(&o.build, "build", "b", false, "build kubemq targets request")
	return cmd
}

func (o *EventsStoreSendOptions) Complete(args []string, transport string) error {
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

func (o *EventsStoreSendOptions) Validate() error {
	return nil
}

func (o *EventsStoreSendOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubemqClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())
	}

	defer func() {
		client.Close()
	}()

	if o.isStream {
		utils.Printlnf("Streaming %d events store messages...", o.messages)
		eventsCh := make(chan *kubemq2.EventStore, 1000)
		eventsResultsCh := make(chan *kubemq2.EventStoreResult, 1000)
		errCh := make(chan error, 10)
		fmt.Println("Sending Stream Events Store:")
		go client.StreamEventsStore(ctx, eventsCh, eventsResultsCh, errCh)
		startTime := time.Now()
		for i := 1; i <= o.messages; i++ {
			msg := client.ES().
				SetChannel(o.channel).
				SetId(uuid.New().String()).
				SetBody([]byte(o.body)).
				SetMetadata(o.metadata)
			printEventStore(msg)
			eventsCh <- msg
			<-eventsResultsCh
		}
		utils.Printlnf("%d events store messages streamed in %s.", o.messages, time.Since(startTime))
		time.Sleep(2 * time.Second)
	} else {
		fmt.Println("Sending Events Store:")
		for i := 1; i <= o.messages; i++ {
			msg := client.ES().
				SetChannel(o.channel).
				SetId(uuid.New().String()).
				SetBody([]byte(o.body)).
				SetMetadata(o.metadata)
			_, err := msg.Send(ctx)
			if err != nil {
				return fmt.Errorf("sending 'events store' body, %s", err.Error())
			}
			printEventStore(msg)
		}

	}
	return nil
}
