package config

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type LicenseOptions struct {
	cfg *config.Config
}

var licenseExamples = `
	# Execute license configuration
	# kubemqctl config license
`
var licenseLong = `Config license command allows to set Kubemqctl license`
var licenseShort = `Config license command allows to set Kubemqctl license`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdLicense(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &LicenseOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "license",
		Aliases: []string{"lic", "l"},
		Short:   licenseShort,
		Long:    licenseLong,
		Example: licenseExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *LicenseOptions) Complete(args []string, cfg *config.Config) error {
	return nil
}

func (o *LicenseOptions) Validate() error {
	return nil
}

func (o *LicenseOptions) Run(ctx context.Context) error {
	err := runLicenseSelection(o.cfg)
	if err != nil {
		return err
	}
	return o.cfg.Save()
}

func runLicenseSelection(cfg *config.Config) error {
	isLicenseData := false
	promptSetLicenseData := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Would you like to set license information?:",
		Default:  false,
	}
	err := survey.AskOne(promptSetLicenseData, &isLicenseData)
	if err != nil {
		return err
	}
	if isLicenseData {
		dataType := ""
		dataTypePrompt := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Select Type:",
			Options: []string{
				"License Key",
				"License Data",
			},
			Default: "License Key",
			Help:    "Select license type input",
		}
		err := survey.AskOne(dataTypePrompt, &dataType)
		if err != nil {
			return err
		}

		switch dataType {
		case "License Key":
			prompt := &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set License Key:",
				Default:  "",
				Help:     "Set License Key:",
			}
			err := survey.AskOne(prompt, &cfg.LicenseKey)
			if err != nil {
				return err
			}
		case "License Data":
			prompt := &survey.Editor{
				Renderer: survey.Renderer{},
				Message:  "Copy & Paste License Data:",
				Default:  "",
				Help:     "Copy & Paste License Data:",
			}
			err := survey.AskOne(prompt, &cfg.LicenseData)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
