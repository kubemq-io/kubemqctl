package operator

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/operator/delete"
	"github.com/kubemq-io/kubemqctl/cmd/operator/install"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/kubemq-io/kubemqctl/cmd/operator/get"

	"github.com/spf13/cobra"
)

var operatorExamples = `
	# Install KubeMQ operator into "kubemq" namespace
	kubemqctl operator install  
	
	# Get KubeMQ operators list 
	kubemqctl operator get  

	# Delete KubeMQ operators 
	kubemqctl operator delete  
`
var operatorLong = `Executes KubeMQ operator management commands`
var operatorShort = `Executes KubeMQ operator management commands`

func NewCmdOperator(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "operator",
		Aliases:   []string{"op"},
		Short:     operatorShort,
		Long:      operatorLong,
		Example:   operatorExamples,
		ValidArgs: []string{"install", "get", "delete"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(install.NewCmdInstall(ctx, cfg))
	cmd.AddCommand(get.NewCmdGet(ctx, cfg))
	cmd.AddCommand(delete.NewCmdDelete(ctx, cfg))

	return cmd
}
