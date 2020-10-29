package manage

import (
	"context"
	"github.com/kubemq-hub/builder/connector/common"
	"github.com/kubemq-hub/builder/manager"
	"github.com/kubemq-io/kubemqctl/cmd/manage/cluster"
	"github.com/kubemq-io/kubemqctl/cmd/manage/connector"
	contextmanager "github.com/kubemq-io/kubemqctl/cmd/manage/context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type Options struct {
	cfg *config.Config
}

var manageExamples = `
	# Execute manage Kubemq components
	kubemqctl manage	
`
var manageLong = `Executes Kubemq manage `
var manageShort = `Executes Kubemq manage command`

func NewCmdManage(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &Options{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "manage",
		Aliases: []string{"m", "mng"},
		Short:   manageShort,
		Long:    manageLong,
		Example: manageExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}
func (o *Options) Complete(args []string) error {
	return nil
}

func (o *Options) Validate() error {
	return nil
}

func (o *Options) Run(ctx context.Context) error {
	contextHandler := contextmanager.NewHandler()
	err := contextHandler.Init(ctx, o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	conHandler := connector.NewHandler()
	err = conHandler.Init(ctx, contextHandler)
	if err != nil {
		return err
	}
	clusterHandler := cluster.NewHandler()
	err = clusterHandler.Init(ctx, contextHandler, conHandler)
	if err != nil {
		return err
	}
	loadedOptions := common.NewDefaultOptions()

	mng := manager.NewManager()
	err = mng.Init(loadedOptions, conHandler, clusterHandler, contextHandler)
	if err := mng.Render(); err != nil {
		return err
	}
	return nil
}
