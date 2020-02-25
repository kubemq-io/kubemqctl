package operator

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	cfg       *config.Config
	namespace string
}

var createExamples = `
	# Install Kubemq operator into "kubemq" namespace
	kubemqctl create operator  

	# Install Kubemq operator into specified namespace
	kubemqctl create operator --namespace my-namespace
 
`
var createLong = `Create command installs kubemq operator into specific namespace`
var createShort = `Create Kubemq Operator`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "operator",
		Aliases: []string{"op", "o"},
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

	cmd.PersistentFlags().StringVarP(&o.namespace, "namespace", "", "kubemq", "namespace name")

	return cmd
}

func (o *CreateOptions) Complete(args []string) error {
	return nil
}

func (o *CreateOptions) Validate() error {
	return nil
}

func (o *CreateOptions) Run(ctx context.Context) error {
	newClient, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}

	dep, err := operatorTypes.CreateDeployment("kubemq-operator", o.namespace)
	if err != nil {
		return err
	}
	_, isUpdate, err := operatorManager.CreateOrUpdateKubemqOperator(dep)
	if err != nil {
		return err
	}
	if isUpdate {
		utils.Printlnf("kubemq operator %s/%s configured.", dep.Namespace, dep.Name)
	} else {
		utils.Printlnf("kubemq operator %s/%s created.", dep.Namespace, dep.Name)
	}
	return nil
}
