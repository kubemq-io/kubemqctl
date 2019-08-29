package commands

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

var commandsExamples = ``
var commandsLong = ``
var commandsShort = `Execute KubeMQ commands based commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCommands(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "commands",
		Aliases: []string{"c", "cmd"},
		Short:   commandsShort,
		Long:    commandsShort,
		Example: commandsExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdCommandsSend(cfg))
	cmd.AddCommand(NewCmdCommandsReceive(cfg))
	cmd.AddCommand(NewCmdCommandsAttach(cfg))

	return cmd
}
