package manage

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/manage/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/manage/connector"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var manageExamples = `
	# Execute manage Kubemq components
	kubemqctl manage	
	
	# Execute manage Kubemq clusters
	kubemqctl manage clusters

	# Execute manage Kubemq connectors
	kubemqctl manage connectors
`
var manageLong = `Executes Kubemq manage command`
var manageShort = `Executes Kubemq manage command`

func NewCmdManage(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "manage",
		Aliases: []string{"m", "mng"},
		Short:   manageShort,
		Long:    manageLong,
		Example: manageExamples,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdManage(ctx, cfg))
	cmd.AddCommand(connector.NewCmdManage(ctx, cfg))
	return cmd
}
