package deploy

import (
	"context"

	"github.com/kubemq-io/kubetools/pkg/config"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type DeployOptions struct {
	cfg *config.Config
}

var deployExamples = `
`
var deployLong = `Stream deploy from pods`
var deployShort = `Stream deploy from pods`

func NewCmdDeploy(cfg *config.Config) *cobra.Command {
	o := &DeployOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "deploy",
		Aliases: []string{"lgs"},
		Short:   deployShort,
		Long:    deployLong,
		Example: deployExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *DeployOptions) Complete(args []string) error {

	return nil
}

func (o *DeployOptions) Validate() error {
	return nil
}

func (o *DeployOptions) Run(ctx context.Context) error {

	return nil
}
