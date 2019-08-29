package events

import (
	"context"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type EventsMetricsOptions struct {
	cfg *config.Config
	top int
}

var eventsMetricsExamples = `
	
`
var eventsMetricsLong = `get events channels metrics`
var eventsMetricsShort = `get events channels metrics`

func NewCmdEventsMetrics(cfg *config.Config) *cobra.Command {
	o := &EventsMetricsOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "metrics",
		Aliases: []string{"a", "att", "at"},
		Short:   eventsMetricsShort,
		Long:    eventsMetricsLong,
		Example: eventsMetricsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Complete(args, cfg.ConnectionType))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	return cmd
}

func (o *EventsMetricsOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *EventsMetricsOptions) Validate() error {
	return nil
}

func (o *EventsMetricsOptions) Run(ctx context.Context) error {

	return nil
}
