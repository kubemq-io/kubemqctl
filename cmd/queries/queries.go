package queries

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var queriesExamples = `
	# Execute send queries command
	# kubemqctl queries send

	# Execute receive queries command
	# kubemqctl queries receive

	# Execute attach to queries command
	# kubemqctl queries attach

`
var queriesLong = `Execute KubeMQ 'queries' RPC based commands`
var queriesShort = `Execute KubeMQ 'queries' RPC based commands`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdQueries(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "queries",
		Aliases: []string{"query", "qry"},
		Short:   queriesShort,
		Long:    queriesLong,
		Example: queriesExamples,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdQueriesSend(cfg))
	cmd.AddCommand(NewCmdQueriesReceive(cfg))
	cmd.AddCommand(NewCmdQueriesAttach(cfg))

	return cmd
}
