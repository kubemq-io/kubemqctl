package events

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/metrics"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

var timeFrames = []string{"1h", "6h", "12h", "1d", "2d", "7d", "30d"}

type EventsMetricsOptions struct {
	cfg *config.Config
	top int
	tf  string
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
		Aliases: []string{"m", "met", "metric"},
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
	cmd.PersistentFlags().IntVarP(&o.top, "top", "t", 10, "set how many intervals to retrieve")
	cmd.PersistentFlags().StringVarP(&o.tf, "time-frame", "d", "1h", `set time frame duration interval - "1h","6h","12h","1d","2d","7d","30d"`)
	return cmd
}

func (o *EventsMetricsOptions) Complete(args []string, transport string) error {
	return nil
}

func (o *EventsMetricsOptions) Validate() error {

	for _, tf := range timeFrames {
		if tf == o.tf {
			return nil
		}
	}
	return fmt.Errorf(`invalid time frame, should be one of"1h","6h","12h","1d","2d","7d","30d"`)
}

func (o *EventsMetricsOptions) Run(ctx context.Context) error {
	err := metrics.PrintMetrics(ctx, os.Stdout, o.cfg, "1", "1", o.tf, o.top)
	if err != nil {
		return err
	}
	return nil
}
