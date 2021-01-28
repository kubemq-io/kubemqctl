package generate

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/generate/authentication"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var generateExamples = `
	# Execute generate authentication commands
 	kubemqctl generate auth

`
var generateLong = `Generate various kubemq related artifacts`
var generateShort = `Generate various kubemq related artifacts`

func NewCmdGenerate(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen"},
		Short:   generateShort,
		Long:    generateLong,
		Example: generateExamples,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(authentication.NewCmdAuthentication(ctx, cfg))
	return cmd
}
