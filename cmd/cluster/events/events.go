package events

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"strings"
)

type EventsOptios struct {
	cfg *config.Config
}

var eventsExamples = `
 	# Show KubeMQ cluster events
	kubemqctl cluster events
`
var eventsLong = `Events command allows to show a real-time KubeMQ cluster events`
var eventsShort = `Show KubeMQ cluster events command`

func NewCmdEvents(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &EventsOptios{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "events",
		Aliases: []string{"e", "ev"},
		Short:   eventsShort,
		Long:    eventsLong,
		Example: eventsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *EventsOptios) Complete(args []string) error {
	return nil
}

func (o *EventsOptios) Validate() error {

	return nil
}

func (o *EventsOptios) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	list, err := c.GetKubeMQClusters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no KubeMQ clusters were found to show events")
	}

	selection := ""
	if len(list) == 1 {
		selection = list[0]
	} else {
		selected := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       "Select KubeMQ cluster to show events",
			Options:       list,
			Default:       list[0],
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
	pair := strings.Split(selection, "/")
	utils.Printlnf("Show real-time events for %s/%s cluster:", pair[0], pair[1])
	go c.PrintEvents(ctx, pair[0], pair[1])
	<-ctx.Done()
	return nil
}
