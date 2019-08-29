package events

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/attach"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type EventsAttachOptions struct {
	cfg       *config.Config
	include   []string
	exclude   []string
	resources []string
}

var eventsAttachExamples = `
	# attach to all events channels and output running messages
	kubetools events attach *
	
	# attach to some-events events channel and output running messages
	kubetools events attach some-events

	# attach to some-events1 and some-events2 events channels and output running messages
	kubetools events attach some-events1 some-events2 

	# attach to some-events events channel and output running messages filter by include regex (some*)
	kubetools events attach some-events -i some*

	# attach to some-events events channel and output running messages filter by exclude regex (not-some*)
	kubetools events attach some-events -e not-some*
`
var eventsAttachLong = `attach to events channels`
var eventsAttachShort = `attach to events channels`

func NewCmdEventsAttach(cfg *config.Config) *cobra.Command {
	o := &EventsAttachOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "attach",
		Aliases: []string{"a", "att", "at"},
		Short:   eventsAttachShort,
		Long:    eventsAttachLong,
		Example: eventsAttachExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Complete(args, cfg.ConnectionType))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringArrayVarP(&o.include, "include", "i", []string{}, "set (regex) strings to include")
	cmd.PersistentFlags().StringArrayVarP(&o.exclude, "exclude", "e", []string{}, "set (regex) strings to exclude")
	return cmd
}

func (o *EventsAttachOptions) Complete(args []string, transport string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing channel argument")

	}

	for _, a := range args {
		rsc := fmt.Sprintf("events/%s", a)
		o.resources = append(o.resources, rsc)
		utils.Printlnf("adding '%s' to attach list", a)
	}
	return nil
}

func (o *EventsAttachOptions) Validate() error {
	return nil
}

func (o *EventsAttachOptions) Run(ctx context.Context) error {
	err := attach.Run(ctx, o.cfg, o.resources, o.include, o.exclude)
	if err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}
