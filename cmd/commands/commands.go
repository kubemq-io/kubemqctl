package commands

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var commandsExamples = `
	# Execute send commands 
	kubemqctl commands send

	# Execute receive commands
	kubemqctl commands receive

	# Execute attach to 'commands' channel
	kubemqctl commands attach
`
var commandsLong = `Execute Kubemq 'commands' RPC commands`
var commandsShort = `Execute Kubemq 'commands' RPC commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCommands(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "commands",
		Aliases:   []string{"cmd"},
		Short:     commandsShort,
		Long:      commandsLong,
		Example:   commandsExamples,
		ValidArgs: []string{"send", "receive", "attach"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdCommandsSend(ctx, cfg))
	cmd.AddCommand(NewCmdCommandsReceive(ctx, cfg))
	cmd.AddCommand(NewCmdCommandsAttach(ctx, cfg))
	return cmd
}
