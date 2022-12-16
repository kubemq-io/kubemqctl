package build

import (
	"context"

	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/pkg/browser"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

func NewCmdBuild(ctx context.Context, cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "Load KubeMQ builder in browser",
		Long:    "Load KubeMQ builder in browser",
		Example: "kubemqctl build",
		Run: func(cmd *cobra.Command, args []string) {
			utils.Println("Loading KubeMQ build in browser")
			utils.CheckErr(browser.OpenURL("https://build.kubemq.io"))
		},
	}

	return cmd
}
