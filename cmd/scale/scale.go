package scale

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/scale/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var scaleExamples = `
	# Execute scale Kubemq cluster
	kubemqctl scale cluster 5	
`
var scaleLong = `Executes Kubemq scale commands`
var scaleShort = `Executes Kubemq scale commands`

func NewCmdScale(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "scale",
		Aliases:   []string{"sc"},
		Short:     scaleShort,
		Long:      scaleLong,
		Example:   scaleExamples,
		ValidArgs: []string{"cluster"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdScale(ctx, cfg))

	return cmd
}
