package scale

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strings"

	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type ScaleOptions struct {
	cfg   *config.Config
	scale int
}

var scaleExamples = `
	# Scale StatufulSet to 5
	kubetools cluster cluster scale 5

	# Scale StatufulSet to 0
	kubetools cluster cluster scale 0
`
var scaleLong = `Scale KubeMQ cluster`
var scaleShort = `Scale KubeMQ cluster`

func NewCmdScale(cfg *config.Config) *cobra.Command {
	o := &ScaleOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "scale",
		Aliases: []string{"scl", "sc"},
		Short:   scaleShort,
		Long:    scaleLong,
		Example: scaleExamples,
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

func (o *ScaleOptions) Complete(args []string) error {

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
	prompt := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Scale for KubeMQ cluster:",
		Options:  list,
		Default:  list[0],
	}
	err = survey.AskOne(prompt, &selection)
	if err != nil {
		return err
	}
	pair := strings.Split(selection, "/")

	promptScale := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set Scale: ",
		Default:  "3",
		Help:     "",
	}
	err = survey.AskOne(promptScale, &o.scale)
	if err != nil {
		return err
	}
	utils.Println("Start scaling:")
	err = c.Scale(ctx, pair[0], pair[1], int32(o.scale))
	if err != nil {
		return err
	}
	return nil
}
