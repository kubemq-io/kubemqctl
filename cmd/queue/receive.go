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
	# receive 1 messages from a queue and wait for 2 seconds (default)
	kubetools queue receive some-channel

	# receive 3 messages from a queue and wait for 5 seconds
	kubetools queue receive some-channel -m 3 -T 5
`
var queueReceiveLong = `receive a message from a queue`
var queueReceiveShort = `receive a message from a queue`

func NewCmdQueueReceive(cfg *config.Config, opts *QueueOptions) *cobra.Command {
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
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	cmd.PersistentFlags().IntVarP(&o.messages, "messages", "m", 1, "set how many messages we want to get from queue")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait-timeout", "T", 2, "set how many seconds to wait for queue messages")
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "set watch on queue channel")

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
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create kubemq client, %s", err.Error())

	}
	defer func() {
		client.Close()
	}()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "CLIENT_ID\tCHANNEL\tID\tMETADATA\tBODY")
	w.Flush()

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
			utils.Println(fmt.Errorf("receive queue message,aa %s", res.Error).Error())
		}

		if res != nil && res.MessagesReceived > 0 {
			printItems(res.Messages)
		}
		if !o.watch {
			return nil
		}
	}

	//for _, item := range res.Messages {
	//
	//	utils.Printlnf("[channel: %s] [client id: %s] -> {id: %s, metadata: %s, body: %s}", item.Channel, item.ClientId, item.Id, item.Metadata, item.Body)
	//}

}

func printItems(items []*kubemq2.QueueMessage) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	for _, item := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", item.ClientId, item.Channel, item.Id, item.Metadata, item.Body)
	}
	w.Flush()
}
