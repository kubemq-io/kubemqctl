package cluster

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"strconv"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ScaleOptions struct {
	cfg   *config.Config
	scale int
}

var scaleExamples = `
	# Scale Kubemq cluster  
	kubemqctl scale cluster 5
`
var scaleLong = `Scale command allows to scale Kubemq cluster replicas`
var scaleShort = `Scale Kubemq cluster replicas command`

func NewCmdScale(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ScaleOptions{
		cfg:   cfg,
		scale: -1,
	}
	cmd := &cobra.Command{

		Use:     "cluster",
		Aliases: []string{"clusters", "c"},
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
	clusterManager, err := cluster.NewManager(c)

	if err != nil {
		return err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}

	if len(clusters.List()) == 0 {
		return fmt.Errorf("no Kubemq cluster to scale")
	}
	selection := ""

	if len(clusters.List()) == 1 {
		selection = clusters.List()[0]
	} else {
		prompt := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Select Kubemq cluster to scale:",
			Options:  clusters.List(),
			Default:  clusters.List()[0],
		}
		err = survey.AskOne(prompt, &selection)
		if err != nil {
			return err
		}
	}

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

	selectedCluster := clusters.Cluster(selection)
	err = clusterManager.ScaleKubemqCluster(selectedCluster, int32(o.scale))
	if err != nil {
		return err
	}
	utils.Printlnf("Scaling cluster %s to %d completed", selection, o.scale)
	return nil
}
