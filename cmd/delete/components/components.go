package components

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	client2 "github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
)

type DeleteOptions struct {
	cfg *config.Config
}

var deleteExamples = `
 	# Delete components
	kubemqctl delete components
`

var (
	deleteLong  = `Delete one or more Kubemq components`
	deleteShort = `Delete Kubemq components`
)

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &DeleteOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "components",
		Aliases: []string{"comp", "cm", "cmp"},
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

func (o *DeleteOptions) runCluster(ctx context.Context, client *client2.Client) error {
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
		return nil
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

func (o *DeleteOptions) runConnector(ctx context.Context, client *client2.Client) error {
	connectorManager, err := connector.NewManager(client)
	if err != nil {
		return err
	}

	operatorManager, err := operator.NewManager(client)
	if err != nil {
		return err
	}
	connectors, err := connectorManager.GetKubemqConnectors()
	if err != nil {
		return err
	}
	if len(connectors.List()) == 0 {
		return nil
	}

	selection := []string{}
	multiSelected := &survey.MultiSelect{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq connectors to delete",
		Options:       connectors.List(),
		Default:       nil,
		Help:          "Select Kubemq connectors to delete",
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
		Help:     "Confirm Kubemq connector deletion",
	}
	err = survey.AskOne(promptConfirm, &areYouSure)
	if err != nil {
		return err
	}
	if !areYouSure {
		return nil
	}
	for _, selected := range selection {
		connector := connectors.Connector(selected)
		if !operatorManager.IsKubemqOperatorExists(connector.Namespace) {
			operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", connector.Namespace)
			if err != nil {
				return err
			}
			_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
			if err != nil {
				return err
			}
			utils.Printlnf("Kubemq operator %s/kubemq-operator created.", connector.Namespace)
		} else {
			utils.Printlnf("Kubemq operator %s/kubemq-operator exists", connector.Namespace)
		}
		err := connectorManager.DeleteKubemqConnector(connector)
		if err != nil {
			return err
		}
		utils.Printlnf("Kubemq connector %s deleted.", selected)
	}
	return nil
}

func (o *DeleteOptions) runOperator(ctx context.Context, client *client2.Client) error {
	operatorManager, err := operator.NewManager(client)
	if err != nil {
		return err
	}

	operators, err := operatorManager.GetKubemqOperatorsDeployments()
	if err != nil {
		return nil
	}

	if len(operators) > 0 {
		var opList []string
		opMap := map[string]*appsv1.Deployment{}
		for _, deployment := range operators {
			opList = append(opList, deployment.Namespace)
			opMap[deployment.Namespace] = deployment
		}

		selection := []string{}
		multiSelected := &survey.MultiSelect{
			Renderer:      survey.Renderer{},
			Message:       "Select Kubemq operator namespace to delete",
			Options:       opList,
			Default:       nil,
			Help:          "Select Kubemq operator namespace to delete",
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
			Help:     "Confirm Kubemq operator deletion",
		}
		err = survey.AskOne(promptConfirm, &areYouSure)
		if err != nil {
			return err
		}
		if !areYouSure {
			return nil
		}
		for _, name := range selection {
			err = operatorManager.DeleteKubemqOperatorDeployment(opMap[name])
			if err != nil {
				return fmt.Errorf("error delete operator for namespace %s, error: %s", name, err.Error())
			}
			utils.Printlnf("Kubemq Operator at namespace %s deleted.", name)
		}
	}

	return nil
}

func (o *DeleteOptions) Run(ctx context.Context) error {
	client, err := client2.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	utils.Printlnf("Getting Kubemq clusters to delete...")
	if err := o.runCluster(ctx, client); err != nil {
		return err
	}
	utils.Printlnf("Getting Kubemq connectors to delete...")
	if err := o.runConnector(ctx, client); err != nil {
		return err
	}

	utils.Printlnf("Getting Kubemq operators to delete...")
	if err := o.runOperator(ctx, client); err != nil {
		return err
	}

	return nil
}
