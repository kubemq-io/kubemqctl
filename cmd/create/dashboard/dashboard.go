package dashboard

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/dashboard"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"reflect"
)

type CreateOptions struct {
	cfg        *config.Config
	isDryRun   bool
	deployOpts *deployOptions
}

var createExamples = `
	# Create default Kubemq Dashboard
	kubemqctl create dashboard
	
	# Create Kubemq dashboard with options - get all flags
	kubemqctl create dashboard --help
`
var createLong = `Create command allows to deploy a Kubemq Dashboard with configuration options`
var createShort = `Create a Kubemq dashboard command`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "dashboard",
		Aliases: []string{"d", "dash", "dashboards"},
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
	cmd.PersistentFlags().BoolVarP(&o.isDryRun, "dry-run", "", false, "generate dashboard configuration without execute")

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
	dashabordManager, err := dashboard.NewManager(newClient)
	if err != nil {
		return err
	}
	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}
	dep := o.deployOpts.getDashboardDeployment()
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
		_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
		if err != nil {
			return nil
		}
		utils.Printlnf("Kubemq operator %s/%s created.", dep.Namespace, dep.Name)
	}

	dashboard, isUpdate, err := dashabordManager.CreateOrUpdateKubemqDashboard(dep)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq dashboard %s/%s configured.", dashboard.Namespace, dashboard.Name)
	} else {
		utils.Printlnf("kubemq dashboard %s/%s created.", dashboard.Namespace, dashboard.Name)
	}
	utils.Println("run: 'kubemqctl get dashboard' for opening dashboard in default browser (it might take couple seconds to load all components")
	return nil

}

func isDefault(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}
