package cluster

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"reflect"
)

type CreateOptions struct {
	cfg        *config.Config
	exportFile bool
	deployOpts *deployOptions
}

var createExamples = `
	# Create default Kubemq cluster
	kubemqctl create cluster
`
var createLong = `Create command allows to deploy a Kubemq cluster with configuration options`
var createShort = `Create a Kubemq cluster command`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "cluster",
		Aliases: []string{"c"},
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
	cmd.PersistentFlags().BoolVarP(&o.exportFile, "export", "e", false, "generate yaml configuration file output (exporting)")

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
	clusterManager, err := cluster.NewManager(newClient)
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}
	dep := o.deployOpts.getClusterDeployment()
	if o.exportFile {
		utils.Printlnf("export to file %s.yaml completed", dep.Name)
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
		_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return nil
		}
		utils.Printlnf("Kubemq operator %s/%s created.", dep.Namespace, dep.Name)
	}

	cluster, isUpdate, err := clusterManager.CreateOrUpdateKubemqCluster(dep)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq cluster %s/%s configured.", cluster.Namespace, cluster.Name)
	} else {
		utils.Printlnf("kubemq cluster %s/%s created.", cluster.Namespace, cluster.Name)
	}

	return nil

}

func isDefault(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
