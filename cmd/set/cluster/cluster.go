package cluster

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/cmd/set/cluster/proxy"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ClusterOptions struct {
	cfg *config.Config
}

var clusterExamples = `
	# Execute set Kubemq cluster connection default
	kubemqctl set cluster 
	# Execute set Kubemq cluster proxy
	kubemqctl set cluster proxy	
`
var clusterLong = `Executes set cluster commands`
var clusterShort = `Executes set cluster commands`

func NewCmdCluster(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ClusterOptions{cfg: cfg}
	cmd := &cobra.Command{

		Use:       "cluster",
		Aliases:   []string{"clusters", "c"},
		Short:     clusterShort,
		Long:      clusterLong,
		Example:   clusterExamples,
		ValidArgs: []string{"proxy"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.AddCommand(proxy.NewCmdProxy(ctx, cfg))
	return cmd
}

func (o *ClusterOptions) Complete(args []string) error {
	return nil
}

func (o *ClusterOptions) Validate() error {
	return nil
}

func (o *ClusterOptions) Run(ctx context.Context) error {
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
		return fmt.Errorf("no Kubemq clusters were found")
	}
	selection := ""
	multiSelected := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select default Kubemq cluster",
		Options:       clusters.List(),
		Default:       clusters.List()[0],
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(multiSelected, &selection)
	if err != nil {
		return err
	}

	ns, name := client.StringSplit(selection)
	o.cfg.CurrentNamespace = ns
	o.cfg.CurrentStatefulSet = name

	err = o.cfg.Save()
	if err != nil {
		return err
	}
	utils.Printlnf("Default Kubemq cluster set to %s", selection)
	return nil
}
