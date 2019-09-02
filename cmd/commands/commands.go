package commands

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

var commandsExamples = `
	# Execute send commands 
	# kubetools commands send

	# Execute receive commands
	# kubetools commands receive

	# Execute attach to commands channel
	# kubetools commands attach
`
var commandsLong = `Execute KubeMQ 'commands' RPC commands`
var commandsShort = `Execute KubeMQ 'commands' RPC commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCommands(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commands",
		Aliases: []string{"c", "cmd"},
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
