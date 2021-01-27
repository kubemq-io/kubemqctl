package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty"
	"github.com/kubemq-io/kubemqctl/pkg/attach"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type QueueAttachOptions struct {
	cfg       *config.Config
	transport string
	include   []string
	exclude   []string
	resources []string
}

var queueAttachExamples = `
	# Attach to all active 'queues' channels and output running messages
	kubemqctl queue attach all
	
	# Attach to some-queue queue channel and output running messages
	kubemqctl queue attach some-queue

	# Attach to some-queue1 and some-queue2 queue channels and output running messages
	kubemqctl queue attach some-queue1 some-queue2 

	# Attach to some-queue queue channel and output running messages filter by include regex (some*)
	kubemqctl queue attach some-queue -i some*

	# Attach to some-queue queue channel and output running messages filter by exclude regex (not-some*)
	kubemqctl queue attach some-queue -e not-some*
`
var queueAttachLong = `Attach command allows to display 'queues' channel content for debugging proposes`
var queueAttachShort = `Attach to 'queues' channels command`

func NewCmdQueueAttach(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &QueueAttachOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "attach",
		Aliases: []string{"a", "att", "at"},
		Short:   queueAttachShort,
		Long:    queueAttachLong,
		Example: queueAttachExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Complete(args, cfg.ConnectionType), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringArrayVarP(&o.include, "include", "i", []string{}, "aet (regex) strings to include")
	cmd.PersistentFlags().StringArrayVarP(&o.exclude, "exclude", "e", []string{}, "set (regex) strings to exclude")
	return cmd
}

func (o *QueueAttachOptions) Complete(args []string, transport string) error {
	o.transport = transport
	if len(args) == 0 {
		return fmt.Errorf("missing channel argument")

	}
	if len(args) == 1 && args[0] == "all" {
		utils.Println("retrieve all active 'queues' channels list...")
		resp := &Response{}
		queues := &Queues{}

		r, err := resty.R().SetResult(resp).SetError(resp).Get(fmt.Sprintf("%s/v1/stats/queues", o.cfg.GetApiHttpURI()))
		if err != nil {
			return err
		}
		if !r.IsSuccess() {
			return fmt.Errorf("not available in current Kubemq version, consider upgrade Kubemq version")
		}
		if resp.Error {
			return fmt.Errorf(resp.ErrorString)
		}
		err = json.Unmarshal(resp.Data, queues)
		if err != nil {
			return err
		}
		utils.Printlnf("found %d active 'queues' channels.", queues.Total)
		for _, q := range queues.Queues {
			rsc := fmt.Sprintf("queue/%s", q.Name)
			o.resources = append(o.resources, rsc)
			utils.Printlnf("adding '%s' to attach list", q.Name)
		}
		return nil
	}
	for _, a := range args {
		rsc := fmt.Sprintf("queue/%s", a)
		o.resources = append(o.resources, rsc)
		utils.Printlnf("adding '%s' to attach list", a)
	}
	return nil
}

func (o *QueueAttachOptions) Validate() error {
	return nil
}

func (o *QueueAttachOptions) Run(ctx context.Context) error {
	err := attach.Run(ctx, o.cfg, o.resources, o.include, o.exclude)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
