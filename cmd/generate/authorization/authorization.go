package authorization

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var authorizationExamples = `
	# Execute generate authorization policy file
 	kubemqctl generate az policy
`
var authorizationLong = `Generate and verify Authentication certificates and tokens`
var authorizationShort = `Generate and verify Authentication certificates and tokens`

func NewCmdAuthorization(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authorization",
		Aliases: []string{"access", "acc", "az"},
		Short:   authorizationShort,
		Long:    authorizationLong,
		Example: authorizationExamples,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdPolicy(ctx, cfg))
	return cmd
}
