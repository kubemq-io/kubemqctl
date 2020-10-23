package connector

import (
	"context"
	"github.com/kubemq-hub/builder/manager"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ManageOptions struct {
	cfg *config.Config
}

var manageExamples = `
	# Manage Kubemq connectors
	kubemqctl manage connectors
`
var manageLong = `Manage command allows to manage Kubemq connectors`
var manageShort = `Manage a Kubemq connectors command`

func NewCmdManage(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ManageOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "connectors",
		Aliases: []string{"con", "connector"},
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

func (o *ManageOptions) Complete(args []string) error {
	return nil
}

func (o *ManageOptions) Validate() error {
	return nil
}

func (o *ManageOptions) Run(ctx context.Context) error {
	h := NewHandler()
	err := h.Init(ctx, o.cfg)
	if err != nil {
		return err
	}
	mng := manager.NewConnectorsManager(h)
	if err := mng.Render(); err != nil {
		return err
	}
	return nil
}
