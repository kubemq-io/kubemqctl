package build

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/web"
	"github.com/pkg/browser"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubemq-io/kubemqctl/pkg/utils"

	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/spf13/cobra"
)

type BuildOptions struct {
	cfg *config.Config
}

func NewCmdBuild(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &BuildOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "Load KubeMQ builder in browser",
		Long:    "Load KubeMQ builder in browser",
		Example: "kubemqctl build",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))

		},
	}

	return cmd
}

func (o *BuildOptions) Complete(args []string) error {
	if o.cfg.WebPort <= 0 {
		o.cfg.WebPort = 55000
	}
	return nil
}

func (o *BuildOptions) Validate() error {
	return nil
}

func (o *BuildOptions) Run(ctx context.Context) error {
	server := web.NewServer()
	if err := server.Init(ctx, o.cfg); err != nil {
		return err
	}
	if err := server.Start(); err != nil {
		return fmt.Errorf("error starting web server, %s", err.Error())
	}
	utils.Printlnf("Loading KubeMQ builder in browser on port %d", o.cfg.WebPort)
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)

	utils.CheckErr(browser.OpenURL(fmt.Sprintf("http://localhost:%d", o.cfg.WebPort)))
	<-gracefulShutdown
	utils.Println("Shutting down KubeMQ builder")
	server.Stop()
	return nil
}
