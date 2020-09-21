package connector

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	cfg        *config.Config
	isDryRun   bool
	deployOpts *deployOptions
}

var createExamples = `
	# Create Kubemq connector
	kubemqctl create connector
	
	# Create Kubemq connector with options - get all flags
	kubemqctl create connector --help
`
var createLong = `Create command allows to deploy a Kubemq connector with configuration options`
var createShort = `Create a Kubemq connector command`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "connector",
		Aliases: []string{"cn", "con"},
		Short:   createShort,
		Long:    createLong,
		Example: createExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	o.deployOpts = defaultDeployOptions(cmd)
	cmd.PersistentFlags().BoolVarP(&o.isDryRun, "dry-run", "", false, "generate connector configuration without execute")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {

	if err := o.deployOpts.complete(); err != nil {
		return err
	}

	return nil
}

func (o *CreateOptions) Validate() error {
	if err := o.deployOpts.validate(); err != nil {
		return err
	}
	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	newClient, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	connectorManager, err := connector.NewManager(newClient)
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}
	dep := o.deployOpts.getConnectorDeployment()
	if o.isDryRun {
		utils.PrintlnfNoTitle(dep.String())
		return nil
	}
	op, err := operatorManager.GetKubemqOperator("kubemq-operator", dep.Namespace)
	if err != nil {
		return err
	}
	if err := op.IsValid(); err != nil {
		operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", dep.Namespace)
		if err != nil {
			return err
		}
		_, isUpdated, err := operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return nil
		}
		if isUpdated {
			utils.Printlnf("Kubemq operator %s/kubemq-operator configured.", dep.Namespace)
		} else {
			utils.Printlnf("Kubemq operator %s/kubemq-operator created.", dep.Namespace)
		}

	}

	connector, isUpdate, err := connectorManager.CreateOrUpdateKubemqConnector(dep)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq connector %s/%s configured.", connector.Namespace, connector.Name)
	} else {
		utils.Printlnf("kubemq connector %s/%s created.", connector.Namespace, connector.Name)
	}

	return nil

}
