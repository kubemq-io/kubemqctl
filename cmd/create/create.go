package create

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/create/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/create/operator"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/spf13/cobra"
)

var createExamples = `
	# Execute create Kubemq cluster
	kubemqctl create cluster	
	
	# Execute create Kubemq operator
	kubemqctl create operator
`
var createLong = `Executes Kubemq create commands`
var createShort = `Executes Kubemq create commands`

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "create",
		Aliases:   []string{"c"},
		Short:     createShort,
		Long:      createLong,
		Example:   createExamples,
		ValidArgs: []string{"cluster", "operator"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdCreate(ctx, cfg))
	cmd.AddCommand(operator.NewCmdCreate(ctx, cfg))

	return cmd
}
