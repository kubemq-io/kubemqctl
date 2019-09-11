package queue

import (
	"context"
	"fmt"
	kubemq2 "github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
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
	# Receive 1 messages from a queue and wait for 10 seconds (default)
	kubetools queue receive some-channel

	# Receive 3 messages from a queue and wait for 5 seconds
	kubetools queue receive some-channel -m 3 -T 5

	# Watching queue channel messages
	kubetools queue receive some-channel -w
`
var queueReceiveLong = `Receive a messages from a queue channel`
var queueReceiveShort = `Receive a messages from a queue channel`

func NewCmdQueueReceive(cfg *config.Config) *cobra.Command {
	o := &QueueReceiveOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "receive",
		Aliases: []string{"r", "rec"},
		Short:   queueReceiveShort,
		Long:    queueReceiveLong,
		Example: queueReceiveExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "Set how many messages we want to get from queue")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait-timeout", "T", 10, "Set how many seconds to wait for queue messages")
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "Set watch on queue channel")

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
		utils.Printlnf("Watching %s queue channel, waiting for messages...", o.channel)
	} else {
		utils.Printlnf("Pulling messages from %s queue channel, waiting for %d seconds...", o.channel, o.wait)
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
			utils.Println(fmt.Errorf("receive queue messagess, %s", err.Error()).Error())
		}
		if res.IsError {
			utils.Println(fmt.Errorf("receive queue message %s", res.Error).Error())
		}

		if res != nil && res.MessagesReceived > 0 {
			printItems(res.Messages)
		} else {
			if !o.watch {
				utils.Println("No new messages in queue")
			}

		}
		if !o.watch {
			return nil
		}
	}

}

func printItems(items []*kubemq2.QueueMessage) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	for _, item := range items {
		fmt.Fprintf(w, "[channel: %s]\t[id: %s]\t[metadata: %s]\t[body: %s]\n", item.Channel, item.Id, item.Metadata, item.Body)
	}
	w.Flush()
}
