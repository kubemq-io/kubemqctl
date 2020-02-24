package events

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"strings"
)

type EventsOptions struct {
	cfg     *config.Config
	cluster string
}

var eventsExamples = `
 	# Show Kubemq cluster events
	kubemqctl get cluster events
`
var eventsLong = `Events command allows to show a real-time Kubemq cluster events`
var eventsShort = `Show Kubemq cluster events command`

func NewCmdEvents(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &EventsOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "events",
		Aliases: []string{"e", "ev"},
		Short:   eventsShort,
		Long:    eventsLong,
		Example: eventsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringVarP(&o.cluster, "cluster", "c", "", "set cluster namespace/name for events")

	return cmd
}

func (o *EventsOptions) Complete(args []string) error {
	return nil
}

func (o *EventsOptions) Validate() error {

	return nil
}

func (o *EventsOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	pair := []string{"", ""}

	if o.cluster == "" {
		clusterManager, err := cluster.NewManager(c)

		if err != nil {
			return err
		}
		clusters, err := clusterManager.GetKubemqClusters()
		if err != nil {
			return err
		}

		if len(clusters.List()) == 0 {
			return fmt.Errorf("no Kubemq clusters were found to show events")
		}

		selection := ""
		if len(clusters.List()) == 1 {
			selection = clusters.List()[0]
		} else {
			selected := &survey.Select{
				Renderer:      survey.Renderer{},
				Message:       "Select Kubemq cluster to show events",
				Options:       clusters.List(),
				Default:       clusters.List()[0],
				PageSize:      0,
				VimMode:       false,
				FilterMessage: "",
				Filter:        nil,
			}
			err = survey.AskOne(selected, &selection)
			if err != nil {
				return err
			}

		}
		pair = strings.Split(selection, "/")
	} else {
		pair = strings.Split(o.cluster, "/")
	}
	if pair[0] == "" || pair[1] == "" {
		return fmt.Errorf("no cluster namespace/name set for event streaming")
	}
	utils.Printlnf("Show real-time events for %s/%s cluster:", pair[0], pair[1])
	go c.PrintEvents(ctx, pair[0], pair[1])
	<-ctx.Done()
	return nil
}
