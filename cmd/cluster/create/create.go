package create

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	cfg        *config.Config
	exportFile bool
	deployOpts *deployOptions
}

var createExamples = `
	# Create default KubeMQ cluster
	kubemqctl cluster create 

	# Create default KubeMQ cluster and watch events and status
	kubemqctl cluster create -w -s

	# Import KubeMQ cluster yaml file  
	kubemqctl cluster create -f kubemq-cluster.yaml

	# Export KubeMQ cluster yaml file    
	kubemqctl cluster create -e 
`
var createLong = `Create command allows to deploy a KubeMQ cluster with configuration options`
var createShort = `Create a KubeMQ cluster command`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &CreateOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "create",
		Aliases: []string{"c", "cr"},
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
	mng, err := manager.NewManager(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	dep := o.deployOpts.getClusterDeployment()
	utils.Printlnf("Create KubeMQ cluster started...")

	if o.exportFile {
		utils.Printlnf("export to file %s.yaml completed", dep.Name)
		return nil
	}
	err = mng.DeployKubeMQCluster(dep)
	if err != nil {
		return err
	}

	return nil

}
