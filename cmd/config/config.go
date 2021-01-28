package config

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ConfigOptions struct {
	Cfg *config.Config
}

var configExamples = `
	# Run Kubemqctl configuration wizard
	# kubemqctl config
`
var configLong = `Config command allows to set Kubemqctl configuration with a wizard`
var configShort = `Run Kubemqctl configuration wizard command`

func NewCmdConfig(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ConfigOptions{
		Cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf", "con"},
		Short:   configShort,
		Long:    configLong,
		Example: configExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.AddCommand(NewCmdConnection(ctx, cfg))
	cmd.AddCommand(NewCmdContext(ctx, cfg))
	cmd.AddCommand(NewCmdAccess(ctx, cfg))
	cmd.AddCommand(NewCmdLicense(ctx, cfg))
	return cmd
}

func (o *ConfigOptions) Complete(args []string, cfg *config.Config) error {
	o.Cfg = cfg
	return nil
}

func (o *ConfigOptions) Validate() error {
	return nil
}

func (o *ConfigOptions) Run(ctx context.Context) error {
	if err := runConnectionSelection(o.Cfg); err != nil {
		return err
	}
	if o.Cfg.AutoIntegrated {
		if err := runContextSelection(o.Cfg); err != nil {
			return err
		}
	}
	if err := runAccessSelection(o.Cfg); err != nil {
		return err
	}
	if err := runLicenseSelection(o.Cfg); err != nil {
		return err
	}
	err := o.Cfg.Save()
	if err != nil {
		return err
	}
	utils.Println("configuration completed and saved.")
	return nil
}
