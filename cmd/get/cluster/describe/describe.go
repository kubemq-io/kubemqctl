package describe

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DescribeOptions struct {
	cfg *config.Config
}

var describeExamples = `
	# Describe Kubemq cluster to console
	kubemqctl get cluster describe
`
var describeLong = `Describe command allows describing a Kubemq cluster to console`
var describeShort = `Describe Kubemq cluster command`

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
	clusterManager, err := cluster.NewManager(c)
	if err != nil {
		return err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}

	if len(clusters.List()) == 0 {
		return fmt.Errorf("no Kubemq clusters were found to describe")
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

	spec := clusters.Cluster(selection)
	utils.PrintlnfNoTitle(spec.String())
	return nil
}
