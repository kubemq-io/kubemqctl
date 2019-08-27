package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/kubemq"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type QueueStreamOptions struct {
	cfg        *config.Config
	transport  string
	channel    string
	visibility int
	wait       int
}

var queueStreamExamples = `
	# stream queue message in transaction mode
	kubetools queue stream some-channel

	# stream queue message in transaction mode with visibility set to 120 seconds and wait time of 180 seconds
	kubetools queue stream some-channel -v 120 -w 180
`
var queueStreamLong = `receive a message from a queue`
var queueStreamShort = `receive a message from a queue`

func NewCmdQueueStream(cfg *config.Config, opts *QueueOptions) *cobra.Command {
	o := &QueueStreamOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "stream",
		Aliases: []string{"st"},
		Short:   queueStreamShort,
		Long:    queueStreamLong,
		Example: queueStreamExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, opts.transport))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().IntVarP(&o.visibility, "visibility", "v", 30, "set initial visibility seconds")
	cmd.PersistentFlags().IntVarP(&o.wait, "wait", "w", 60, "set how many seconds to wait for queue messages")

	return cmd
}

func (o *QueueStreamOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) >= 1 {
		o.channel = args[0]
		return nil
	}
	return fmt.Errorf("missing channel argument")
}

func (o *QueueStreamOptions) Validate() error {
	return nil
}

func (o *QueueStreamOptions) Run(ctx context.Context) error {
	client, err := kubemq.GetKubeMQClient(ctx, o.transport, o.cfg)
	if err != nil {
		return fmt.Errorf("create send client, %s", err.Error())

	}
	defer utils.CheckErr(client.Close())
	stream := client.NewStreamQueueMessage().SetChannel(o.channel)
	msg, err := stream.Next(ctx, int32(o.visibility), int32(o.wait))
	if err != nil {
		return err
	}
	utils.Printlnf("message received - id: %s, metadata: %s, body: %s", msg.Id, msg.Metadata, string(msg.Body))
	return nil
}
func (o *QueueStreamOptions) prompt() (string, interface{}) {

}
