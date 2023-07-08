package config

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type AccessOptions struct {
	cfg *config.Config
}

var accessExamples = `
	# Execute access configuration
	# kubemqctl config access
`
var accessLong = `Config access command allows to set Kubemqctl access`
var accessShort = `Config access command allows to set Kubemqctl access`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdAccess(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &AccessOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "access",
		Aliases: []string{"acc", "a"},
		Short:   accessShort,
		Long:    accessLong,
		Example: accessExamples,
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

func (o *AccessOptions) Complete(args []string, cfg *config.Config) error {
	return nil
}

func (o *AccessOptions) Validate() error {
	return nil
}

func (o *AccessOptions) Run(ctx context.Context) error {
	err := runAccessSelection(o.cfg)
	if err != nil {
		return err
	}
	return o.cfg.Save()
}

func runAccessSelection(cfg *config.Config) error {
	isConfig := false
	promptIsSecured := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Configure access control",
		Default:  false,
		Help:     "Configure access control",
	}
	err := survey.AskOne(promptIsSecured, &isConfig)
	if err != nil {
		return err
	}

	if isConfig {
		promptIsSecured := &survey.Confirm{
			Renderer: survey.Renderer{},
			Message:  "Set SSL Secured connection ?",
			Default:  false,
			Help:     "Set Kubemq secured connection",
		}
		err := survey.AskOne(promptIsSecured, &cfg.IsSecured)
		if err != nil {
			return err
		}
		if cfg.IsSecured {
			promptCertFile := &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set cert file path:",
				Default:  "",
				Help:     "Set Kubemq cert file path",
			}
			err = survey.AskOne(promptCertFile, &cfg.CertFile)
			if err != nil {
				return err
			}
			_, err = ioutil.ReadFile(cfg.CertFile)
			if err != nil {
				return fmt.Errorf("error loading cert file: %s", err.Error())
			}
		}

		isSetClientId := false
		promptSetClientIDSecured := &survey.Confirm{
			Renderer: survey.Renderer{},
			Message:  "Would you like to set connection ClientId ?",
			Default:  false,
		}
		err = survey.AskOne(promptSetClientIDSecured, &isSetClientId)
		if err != nil {
			return err
		}
		if isSetClientId {
			promptClientId := &survey.Input{
				Renderer: survey.Renderer{},
				Message:  "Set ClientId:",
				Default:  cfg.ClientId,
				Help:     "Set ClientId for every connection",
			}
			err = survey.AskOne(promptClientId, &cfg.ClientId)
			if err != nil {
				return err
			}
		} else {
			cfg.ClientId = ""
		}
		isAuthToken := false
		promptSetAuthToken := &survey.Confirm{
			Renderer: survey.Renderer{},
			Message:  "Would you like to set JWT Authentication token ?",
			Default:  false,
		}
		err = survey.AskOne(promptSetAuthToken, &isAuthToken)
		if err != nil {
			return err
		}
		if isAuthToken {
			promptAuthToken := &survey.Editor{
				Renderer: survey.Renderer{},
				Message:  "Set JWT Authentication token file",
				Default:  cfg.AuthTokenFile,
				Help:     "Set JWT Authentication token file",
			}
			err = survey.AskOne(promptAuthToken, &cfg.AuthTokenFile)
			if err != nil {
				return err
			}
		} else {
			cfg.AuthTokenFile = ""
		}
	}

	return nil
}
