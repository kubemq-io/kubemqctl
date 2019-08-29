package queries

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

var queriesExamples = ``
var queriesLong = ``
var queriesShort = `Execute KubeMQ queries based commands`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdQueries(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "queries",
		Aliases: []string{"query", "qry"},
		Short:   queriesShort,
		Long:    queriesShort,
		Example: queriesExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdQueriesSend(cfg))
	cmd.AddCommand(NewCmdQueriesReceive(cfg))
	cmd.AddCommand(NewCmdQueriesAttach(cfg))

	return cmd
}
