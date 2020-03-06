package get

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/get/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/get/dashboard"
	"github.com/kubemq-io/kubemqctl/cmd/get/operator"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/spf13/cobra"
)

var getExamples = `
	# Execute get Kubemq clusters
	kubemqctl get clusters	
	
	# Execute get Kubemq operators
	kubemqctl get operators	

	# Execute get Kubemq dashboards
	kubemqctl get dashboards
`
var getLong = `Executes Kubemq get commands`
var getShort = `Executes Kubemq get commands`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "get",
		Aliases:   []string{"g"},
		Short:     getShort,
		Long:      getLong,
		Example:   getExamples,
		ValidArgs: []string{"cluster", "operator"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdGet(ctx, cfg))
	cmd.AddCommand(operator.NewCmdGet(ctx, cfg))
	cmd.AddCommand(dashboard.NewCmdGet(ctx, cfg))

	return cmd
}
