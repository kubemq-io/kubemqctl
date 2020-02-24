package set

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/set/cluster"
	c "github.com/kubemq-io/kubemqctl/cmd/set/context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var setExamples = `
	# Execute set Kubemq cluster proxy
	kubemqctl set cluster proxy
	# Execute set kubernetes context
	kubemqctl set context
`
var setLong = `Executes set commands`
var setShort = `Executes set commands`

func NewCmdSet(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "set",
		Short:     setShort,
		Long:      setLong,
		Example:   setExamples,
		ValidArgs: []string{"cluster", "context"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdCluster(ctx, cfg))
	cmd.AddCommand(c.NewCmdContext(ctx, cfg))
	return cmd
}
