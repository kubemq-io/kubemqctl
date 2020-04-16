package queries

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var queriesExamples = `
	# Execute send 'queries'  command
	kubemqctl queries send

	# Execute receive 'queries'  command
	kubemqctl queries receive

	# Execute attach to 'queries' command
	kubemqctl queries attach

`
var queriesLong = `Execute Kubemq 'queries' RPC based commands`
var queriesShort = `Execute Kubemq 'queries' RPC based commands`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdQueries(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:       "queries",
		Aliases:   []string{"query", "qry"},
		Short:     queriesShort,
		Long:      queriesLong,
		Example:   queriesExamples,
		ValidArgs: []string{"send", "receive", "attach"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdQueriesSend(ctx, cfg))
	cmd.AddCommand(NewCmdQueriesReceive(ctx, cfg))
	cmd.AddCommand(NewCmdQueriesAttach(ctx, cfg))

	return cmd
}
