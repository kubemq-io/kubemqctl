package cluster

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	client2 "github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
 	# Delete Kubemq cluster
	kubemqctl delete cluster
`
var deleteLong = `Delete one or more Kubemq clusters`
var deleteShort = `Delete Kubemq cluster`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "cluster",
		Aliases: []string{"c", "clusters"},
		Short:   deleteShort,
		Long:    deleteLong,
		Example: deleteExamples,
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

func (o *DeleteOptions) Complete(args []string) error {
	return nil
}

func (o *DeleteOptions) Validate() error {

	return nil
}

func (o *DeleteOptions) Run(ctx context.Context) error {
	client, err := client2.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	clusterManager, err := cluster.NewManager(client)
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(client)
	if err != nil {
		return err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}
	if len(clusters.List()) == 0 {
		return fmt.Errorf("no Kubemq clusters were found to delete")
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq clusters to delete",
		Options:       clusters.List(),
		Default:       nil,
		Help:          "Select Kubemq clusters to delete",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(multiSelected, &selection)
	if err != nil {
		return err
	}

	areYouSure := false
	promptConfirm := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Are you sure ?",
		Default:  false,
		Help:     "Confirm Kubemq cluster deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, selected := range selection {
		cluster := clusters.Cluster(selected)
		if !operatorManager.IsKubemqOperatorExists(cluster.Namespace) {
			operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", cluster.Namespace)
			if err != nil {
				return err
			}
			_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
			if err != nil {
				return err
			}
			utils.Printlnf("Kubemq operator %s/kubemq-operator created.", cluster.Namespace)
		} else {
			utils.Printlnf("Kubemq operator %s/kubemq-operator exists", cluster.Namespace)
		}
		err := clusterManager.DeleteKubemqCluster(cluster)
		if err != nil {
			return err
		}
		utils.Printlnf("Kubemq cluster %s deleted.", selected)
	}
	return nil
}
