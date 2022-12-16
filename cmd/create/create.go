package create

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

func NewCmdCreate(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},

		Run: func(cmd *cobra.Command, args []string) {
			utils.Println("kubemq create command is deprecated, please use kubemqctl install command")
		},
	}

	return cmd
}
