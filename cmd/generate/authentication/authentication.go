package authentication

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

var authenticationExamples = `
	# Execute generate authentication rsa certificates
 	kubemqctl generate auth certs

`
var authenticationLong = `Generate and verify Authentication certificates and tokens`
var authenticationShort = `Generate and verify Authentication certificates and tokens`

func NewCmdAuthentication(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "authentication",
		Aliases: []string{"auth"},
		Short:   authenticationShort,
		Long:    authenticationLong,
		Example: authenticationExamples,
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(NewCmdCerts(ctx, cfg))
	cmd.AddCommand(NewCmdToken(ctx, cfg))
	return cmd
}
