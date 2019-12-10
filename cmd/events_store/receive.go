package events_store

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	kubemq2 "github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"text/tabwriter"
	"time"
)

type EventsStoreReceiveOptions struct {
	cfg        *config.Config
	transport  string
	channel    string
	group      string
	subOptions kubemq2.SubscriptionOption
}

var eventsReceiveExamples = `
	# Receive messages from an 'events store' channel (blocks until next message)
	kubemqctl events_store receive some-channel

	# Receive messages from an 'events channel' with group(blocks until next message)
	kubemqctl events_store receive some-channel -g G1
`
var eventsReceiveLong = `Receive (Subscribe) command allows to consume messages from an 'events store' with options to set offset parameters`
var eventsReceiveShort = `Receive a messages from an 'events store'`

func NewCmdEventsStoreReceive(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &EventsStoreReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec"},
		Short:   eventsReceiveShort,
		Long:    eventsReceiveLong,
		Example: eventsReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().StringVarP(&o.group, "group", "g", "", "set 'events_store' channel consumer group (load balancing)")
	return cmd
}

func (o *EventsStoreReceiveOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *EventsStoreReceiveOptions) Validate() error {
	return nil
}

func (o *EventsStoreReceiveOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	err = o.promptOptions()
	if err != nil {
		return err
	}

	errChan := make(chan error, 1)
	eventsChan, err := client.SubscribeToEventsStore(ctx, o.channel, o.group, errChan, o.subOptions)

	if err != nil {
		utils.Println(fmt.Errorf("receive 'events store' messages, %s", err.Error()).Error())
	}
	utils.Println("waiting for 'events store' messages...")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.TabIndent)
	for {
		select {
		case ev, opened := <-eventsChan:
			if !opened {
				utils.Println("server disconnected")
				return nil
			}
			fmt.Fprintf(w, "[channel: %s]\t[seq: %d]\t[time: %s(UTC)]\t[id: %s]\t[metadata: %s]\t[body: %s]\n", ev.Channel, ev.Sequence, ev.Timestamp.UTC().Format("2006-01-02 15:04:05"), ev.Id, ev.Metadata, ev.Body)
			w.Flush()
		case err := <-errChan:
			return fmt.Errorf("server disconnected with error: %s", err.Error())
		case <-ctx.Done():
			return nil
		}
	}

}

func (o *EventsStoreReceiveOptions) promptOptions() error {
	action := ""
	prompt := &survey.Select{
		Message: "Start receive events store messages options:",
		Options: []string{"start from new messages only", "start from first message", "start from last message", "start from sequence", "start from time", "start from duration"},
	}
	err := survey.AskOne(prompt, &action)
	if err != nil {
		return err
	}
	switch action {
	case "start from new messages only":
		o.subOptions = kubemq2.StartFromNewEvents()
		return nil
	case "start from first message":
		o.subOptions = kubemq2.StartFromFirstEvent()
		return nil
	case "start from last message":
		o.subOptions = kubemq2.StartFromLastEvent()
		return nil
	case "start from sequence":
		seqStr := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set sequence:",
			Default:  "1",
			Help:     "1 is the first message",
		}

		err := survey.AskOne(prompt, &seqStr)
		if err != nil {
			return err
		}

		seq, err := strconv.Atoi(seqStr)
		if err != nil {
			return err
		}
		o.subOptions = kubemq2.StartFromSequence(seq)
		return nil
	case "start from time":
		timeStr := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set time (UTC):",
			Default:  time.Now().UTC().Add(-1 * time.Minute).Format("2006-01-02 15:04:05"),
			Help:     "Time format '2006-01-02 15:04:05'",
		}

		err := survey.AskOne(prompt, &timeStr)
		if err != nil {
			return err
		}
		t, err := time.Parse("2006-01-02 15:04:05", timeStr)

		if err != nil {
			return err
		}
		o.subOptions = kubemq2.StartFromTime(t.UTC())
		return nil
	case "start from duration":
		durationStr := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  fmt.Sprintf("Set duration (Current UTC time: %s):", time.Now().UTC().Format("2006-01-02 15:04:05")),
			Default:  "1h",
			Help:     "Duration time such 1s, 1h, 24h",
		}

		err := survey.AskOne(prompt, &durationStr)
		if err != nil {
			return err
		}
		d, err := time.ParseDuration(durationStr)

		if err != nil {
			return err
		}
		o.subOptions = kubemq2.StartFromTimeDelta(d)
		return nil
	}
	return fmt.Errorf("invalid input")
}
