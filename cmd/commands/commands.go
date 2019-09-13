package commands

import (
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

var commandsExamples = `
	# Execute send commands 
	# kubemqctl commands send

	# Execute receive commands
	# kubemqctl commands receive

	# Execute attach to commands channel
	# kubemqctl commands attach
`
var commandsLong = `Execute KubeMQ 'commands' RPC commands`
var commandsShort = `Execute KubeMQ 'commands' RPC commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCommands(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commands",
		Aliases: []string{"cmd"},
		Short:   commandsShort,
		Long:    commandsLong,
		Example: commandsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(NewCmdCommandsSend(cfg))
	cmd.AddCommand(NewCmdCommandsReceive(cfg))
	cmd.AddCommand(NewCmdCommandsAttach(cfg))
	return cmd
}
