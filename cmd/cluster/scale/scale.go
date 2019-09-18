package scale

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strconv"
	"strings"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"

	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ScaleOptions struct {
	cfg    *config.Config
	scale  int
	watch  bool
	status bool
}

var scaleExamples = `
	# Scale KubeMQ cluster StatefulSet 
	kubemqctl cluster cluster scale 5

	# Scale KubeMQ cluster StatefulSet with streaming real-time events and status
	kubemqctl cluster scale -w -s 

	# Scale KubeMQ cluster StatefulSet to 0
	kubemqctl cluster scale 0
`
var scaleLong = `Scale command allows ro scale KubeMQ cluster replicas`
var scaleShort = `Scale KubeMQ cluster replicas command`

func NewCmdScale(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ScaleOptions{
		cfg:   cfg,
		scale: -1,
	}
	cmd := &cobra.Command{

		Use:     "scale",
		Aliases: []string{"scl", "sc"},
		Short:   scaleShort,
		Long:    scaleLong,
		Example: scaleExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "watch and print Scale StatefulSet events")
	cmd.PersistentFlags().BoolVarP(&o.status, "status", "s", false, "watch and print Scale StatefulSet status")

	return cmd
}

func (o *ScaleOptions) Complete(args []string) error {
	var err error
	if len(args) > 0 {
		o.scale, err = strconv.Atoi(args[0])
		if err != nil {
			o.scale = -1
		}
	}
	return nil
}

func (o *ScaleOptions) Validate() error {
	return nil
}

func (o *ScaleOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	list, err := c.GetKubeMQClusters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no KubeMQ cluster to scale")
	}
	selection := ""

	if len(list) == 1 {
		selection = list[0]
	} else {
		prompt := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Select KubeMQ cluster to scale:",
			Options:  list,
			Default:  list[0],
		}
		err = survey.AskOne(prompt, &selection)
		if err != nil {
			return err
		}
	}

	pair := strings.Split(selection, "/")

	if o.scale < 0 {
		promptScale := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set Scale: ",
			Default:  "",
			Help:     "",
		}
		err = survey.AskOne(promptScale, &o.scale)
		if err != nil {
			return err
		}
	}

	utils.Println("Scaling started")
	err = c.Scale(ctx, pair[0], pair[1], int32(o.scale))
	if err != nil {
		return err
	}

	if o.watch {
		go c.PrintEvents(ctx, pair[0], pair[1])
	}

	if o.status {
		go c.PrintStatefulSetStatus(ctx, int32(o.scale), pair[0], pair[1])
	}
	if o.status || o.watch {
		<-ctx.Done()

	}

	return nil
}
