package install

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type InstallOptions struct {
	cfg       *config.Config
	namespace string
}

var installExamples = `
	# Install KubeMQ operator into "kubemq" namespace
	kubemqctl operator install  

	# Install KubeMQ operator into specified namespace
	kubemqctl operator install  --namespace my-namespace
 
`
var installLong = `Install command installs kubemq operator into specific namespace`
var installShort = `Install KubeMQ Operator`

func NewCmdInstall(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &InstallOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "install",
		Aliases: []string{"i", "ins"},
		Short:   installShort,
		Long:    installLong,
		Example: installExamples,
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

func (o *InstallOptions) Complete(args []string) error {
	return nil
}

func (o *InstallOptions) Validate() error {
	return nil
}

func (o *InstallOptions) Run(ctx context.Context) error {
	mng, err := manager.NewManager(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	err = mng.DeployOperator(o.namespace)
	if err != nil {
		return err
	}
	utils.Printlnf("KubeMQ operator deployed successfully to %s namespace", o.namespace)
	return nil
}
