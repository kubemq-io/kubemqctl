package events

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type EventsReceiveOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	group     string
}

var eventsReceiveExamples = `
	# Receive messages from an events channel (blocks until next message)
	kubemqctl events receive some-channel

	# Receive messages from an events channel with group (blocks until next message)
	kubemqctl events receive some-channel -g G1

`
var eventsReceiveLong = `Receive a message from events channel`
var eventsReceiveShort = `Receive a message from events channel`

func NewCmdEventsReceive(cfg *config.Config) *cobra.Command {
	o := &EventsReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec"},
		Short:   eventsReceiveShort,
		Long:    eventsReceiveLong,
		Example: eventsReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.group, "group", "g", "", "Set group")
	return cmd
}

func (o *EventsReceiveOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *EventsReceiveOptions) Validate() error {
	return nil
}

func (o *EventsReceiveOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)

	errChan := make(chan error, 1)
	eventsChan, err := client.SubscribeToEvents(ctx, o.channel, o.group, errChan)

	if err != nil {
		utils.Println(fmt.Errorf("receive events messages, %s", err.Error()).Error())
	}
	utils.Println("waiting for events messages...")
	for {
		select {
		case ev, opened := <-eventsChan:
			if !opened {
				utils.Println("server disconnected")
				return nil
			}
			fmt.Fprintf(w, "[channel: %s]\t[id: %s]\t[metadata: %s]\t[body: %s]\n", ev.Channel, ev.Id, ev.Metadata, ev.Body)
			w.Flush()
		case <-ctx.Done():
			return nil
		}
	}

}
