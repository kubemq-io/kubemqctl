package describe

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/dashboard"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DescribeOptions struct {
	cfg *config.Config
}

var describeExamples = `
	# Describe Kubemq dashboard to console
	kubemqctl get dashboard describe
`
var describeLong = `Describe command allows describing a Kubemq dashboard to console`
var describeShort = `Describe Kubemq dashboard command`

func NewCmdDescribe(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DescribeOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "describe",
		Aliases: []string{"des", "d"},
		Short:   describeShort,
		Long:    describeLong,
		Example: describeExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *DescribeOptions) Complete(args []string) error {
	return nil
}

func (o *DescribeOptions) Validate() error {

	return nil
}

func (o *DescribeOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	dashboardManager, err := dashboard.NewManager(c)
	if err != nil {
		return err
	}
	dashboards, err := dashboardManager.GetKubemqDashboardes()
	if err != nil {
		return err
	}

	if len(dashboards.List()) == 0 {
		return fmt.Errorf("no Kubemq dashboards were found to describe")
	}

	selection := ""
	if len(dashboards.List()) == 1 {
		selection = dashboards.List()[0]
	} else {
		selected := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       "Select Kubemq dashboard to describe",
			Options:       dashboards.List(),
			Default:       dashboards.List()[0],
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

	spec := dashboards.Dashboard(selection)
	utils.PrintlnfNoTitle(spec.String())
	return nil
}
