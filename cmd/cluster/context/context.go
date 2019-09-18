package context

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"sort"

	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ContextOptions struct {
	cfg *config.Config
}

var contextExamples = `
	# Select kubernetes cluster context
	kubemqctl cluster context

`
var contextLong = `Context command allows to set the default Kubernetes cluster context`
var contextShort = `Select kubernetes cluster context command`

func NewCmdContext(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ContextOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "context",
		Aliases: []string{"cx", "ct", "ctx"},
		Short:   contextShort,
		Long:    contextLong,
		Example: contextExamples,
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

func (o *ContextOptions) Complete(args []string) error {

	return nil
}

func (o *ContextOptions) Validate() error {
	return nil
}

func (o *ContextOptions) Run(ctx context.Context) error {

	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	contextMap, current, err := c.GetConfigContext()
	if err != nil {
		return err
	}
	list := []string{}
	for key, _ := range contextMap {
		list = append(list, key)
	}
	sort.Strings(list)
	contextSelected := ""
	contextSelect := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select kubernetes cluster context",
		Options:       list,
		Default:       current,
		Help:          "Set kubernetes connection context",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(contextSelect, &contextSelected)
	if err != nil {
		return err
	}
	err = c.SwitchContext(contextSelected)
	if err != nil {
		return err
	}
	utils.Printlnf("kubernetes cluster switched to %s", contextSelected)
	return nil
}
