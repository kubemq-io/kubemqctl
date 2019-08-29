package commands

import (
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

type CommandsOptions struct {
	transport string
}

var commandsExamples = ``
var commandsLong = ``
var commandsShort = `Execute KubeMQ commands based commands`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdCommands(cfg *config.Config) *cobra.Command {
	o := &CommandsOptions{
		transport: "grpc",
	}
	cmd := &cobra.Command{
		Use:     "commands",
		Aliases: []string{"c", "cmd"},
		Short:   commandsShort,
		Long:    commandsShort,
		Example: commandsExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(NewCmdCommandsSend(cfg, o))
	cmd.AddCommand(NewCmdCommandsReceive(cfg, o))
	cmd.AddCommand(NewCmdCommandsAttach(cfg, o))

	cmd.PersistentFlags().StringVarP(&o.transport, "transport", "t", "grpc", "set transport type, grpc or rest")
	return cmd
}
