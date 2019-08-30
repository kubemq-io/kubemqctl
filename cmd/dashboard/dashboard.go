package dashboard

import (
	"context"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/kubemq-io/kubetools/web"
	"github.com/spf13/cobra"
)

type DashboardOptions struct {
	cfg *config.Config
}

var dashboardExamples = ``
var dashboardLong = `Execute KubeMQ dashboard commands`
var dashboardShort = `Execute KubeMQ dashboard commands`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdDashboard(cfg *config.Config) *cobra.Command {
	o := DashboardOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "dashboard",
		Aliases: []string{"web", "dash", "d"},
		Short:   dashboardLong,
		Long:    dashboardShort,
		Example: dashboardExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args))
			utils.CheckErr(o.Validate())
			utils.CheckErr(k8s.SetTransport(ctx, cfg))
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *DashboardOptions) Complete(args []string) error {

	return nil
}

func (o *DashboardOptions) Validate() error {
	return nil
}

func (o *DashboardOptions) Run(ctx context.Context) error {
	s := &web.ServerOptions{
		Cfg:  o.cfg,
		Port: 6700,
		Path: "./web/dist",
	}
	//err := s.Run(ctx)
	//if err != nil {
	//	return err
	//}
	go func() {
		err := s.Download(ctx)
		if err != nil {
			utils.CheckErr(err)
		}

	}()
	err := s.Download(ctx)
	if err != nil {
		return err
	}
	<-ctx.Done()

	return nil
}
