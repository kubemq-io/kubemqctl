package generate

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/generate/authentication"
	"github.com/kubemq-io/kubemqctl/cmd/generate/authorization"
	"github.com/kubemq-io/kubemqctl/cmd/generate/routing"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var generateExamples = `
	# Execute generate authentication commands
 	kubemqctl generate auth

	# Execute generate authorization commands
 	kubemqctl generate az

	# Execute generate smart routing file
 	kubemqctl generate routes
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
	cmd.AddCommand(authorization.NewCmdAuthorization(ctx, cfg))
	cmd.AddCommand(routing.NewCmdRouting(ctx, cfg))
	return cmd
}
