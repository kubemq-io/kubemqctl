package delete

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/delete/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/delete/components"
	"github.com/kubemq-io/kubemqctl/cmd/delete/connector"
	"github.com/kubemq-io/kubemqctl/cmd/delete/operator"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/spf13/cobra"
)

var deleteExample = `
	# Execute delete Kubemq cluster
	kubemqctl delete cluster
	
	# Execute delete Kubemq Operator
	kubemqctl delete operator	
	
	# Execute delete Kubemq Dashboard
	kubemqctl delete dashboard	

   # Execute delete Kubemq Connector
	kubemqctl delete connector

   # Execute delete Kubemq Components
	kubemqctl delete connector
`
var createLong = `Executes delete commands`
var createShort = `Executes delete commands`

func NewCmdDelete(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "delete",
		Aliases:   []string{"d", "del"},
		Short:     createShort,
		Long:      createLong,
		Example:   deleteExample,
		ValidArgs: []string{"cluster", "operator", "dashboard", "connector", "component"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(cluster.NewCmdDelete(ctx, cfg))
	cmd.AddCommand(operator.NewCmdDelete(ctx, cfg))
	cmd.AddCommand(connector.NewCmdDelete(ctx, cfg))
	cmd.AddCommand(components.NewCmdDelete(ctx, cfg))
	return cmd
}
