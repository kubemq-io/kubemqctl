package queue

import (
	"context"
	"fmt"
	kubemq2 "github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/kubemq"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

type QueueReceiveOptions struct {
	cfg       *config.Config
	transport string
	channel   string
	messages  int
	wait      int
	watch     bool
}

var queueReceiveExamples = `
	# Receive 1 messages from a queue channel q1 and wait for 2 seconds (default)
	kubemqctl queue receive q1

	# Receive 3 messages from a queue channel and wait for 5 seconds
	kubemqctl queue receive q1 -m 3 -t 5

	# Watching 'queues' channel messages
	kubemqctl queue receive q1 -w
`
var queueReceiveLong = `Receive command allows to receive one or many messages from a queue channel`
var queueReceiveShort = `Receive a messages from a queue channel command`

func NewCmdQueueReceive(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueueReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec", "subscribe", "sub"},
		Short:   queueReceiveShort,
		Long:    queueReceiveLong,
		Example: queueReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set how many messages we want to get from a queue")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait-timeout", "t", 2, "set how many seconds to wait for 'queues' messages")
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "set watch on 'queues' channel")

	return cmd
}

func (o *QueueReceiveOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueueReceiveOptions) Validate() error {
	return nil
}

func (o *QueueReceiveOptions) Run(ctx context.Context) error {
	if o.watch {
		utils.Printlnf("Watching %s 'queues' channel, waiting for messages...", o.channel)
	} else {
		utils.Printlnf("Pulling %d messages from %s 'queues' channel, waiting for %d seconds...", o.messages, o.channel, o.wait)
	}
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	for {
		res, err := client.RQM().
			SetChannel(o.channel).
			SetWaitTimeSeconds(o.wait).
			SetMaxNumberOfMessages(o.messages).
			Send(ctx)
		if err != nil {
			utils.Println(fmt.Errorf("receive 'queues' messages, %s", err.Error()).Error())
		}
		if res.IsError {
			utils.Println(fmt.Errorf("receive 'queues' messages %s", res.Error).Error())
		}

		if res != nil && res.MessagesReceived > 0 {
			printItems(res.Messages)
		} else if !o.watch {
			utils.Println("No new messages in queue")

		}
		if !o.watch {
			return nil
		}
	}

}

func printItems(items []*kubemq2.QueueMessage) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	for _, item := range items {
		fmt.Fprintf(w, "[channel: %s]\t[id: %s]\t[seq: %d]\t[timestamp: %s]\t[metadata: %s]\t[body: %s]\n",
			item.Channel,
			item.Id,
			item.Attributes.Sequence,
			time.Unix(0, item.Attributes.Timestamp).Format("2006-01-02 15:04:05.999"), item.Metadata,
			item.Body)
	}
	w.Flush()
}
