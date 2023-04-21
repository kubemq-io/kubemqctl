package config

import (
	"context"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ConnectionOptions struct {
	cfg *config.Config
}

var connectionExamples = `
	# Execute connection configuration
	# kubemqctl config connection
`

var (
	connectionLong  = `Config connection command allows to set Kubemqctl connection`
	connectionShort = `Config connection command allows to set Kubemqctl connection`
)

func NewCmdConnection(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ConnectionOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "connection",
		Aliases: []string{"con"},
		Short:   connectionShort,
		Long:    connectionLong,
		Example: connectionExamples,
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

func (o *ConnectionOptions) Complete(args []string, cfg *config.Config) error {
	return nil
}

func (o *ConnectionOptions) Validate() error {
	return nil
}

func (o *ConnectionOptions) Run(ctx context.Context) error {
	err := runConnectionSelection(o.cfg)
	if err != nil {
		return err
	}
	return o.cfg.Save()
}

func runConnectionSelection(cfg *config.Config) error {
	integrationType := ""
	integrationSelect := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Select Connection Type:",
		Options: []string{
			"Kubernetes cluster",
			"Direct",
		},
		Default: "Kubernetes cluster",
		Help:    "Select the location of Kubemq server",
	}
	err := survey.AskOne(integrationSelect, &integrationType)
	if err != nil {
		return err
	}
	switch integrationType {
	case "Kubernetes cluster":
		cfg.AutoIntegrated = true
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Select kube config path (press Enter for default):",
			Default:  "",
			Help:     "Set kube.config file path if not kubectl default",
		}
		err := survey.AskOne(prompt, &cfg.KubeConfigPath)
		if err != nil {
			return err
		}
		cfg.GrpcPort = 50000
		cfg.RestPort = 9090
		cfg.ApiPort = 8080
		cfg.Host = "localhost"
	case "Direct":
		cfg.AutoIntegrated = false
		promptHost := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set Kubemq Host:",
			Default:  "localhost",
			Help:     "Set Kubemq host",
		}
		err = survey.AskOne(promptHost, &cfg.Host, survey.WithValidator(survey.MinLength(1)))
		if err != nil {
			return err
		}
		promptGrpcPort := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set Kubemq gRPC port interface:",
			Default:  "50000",
			Help:     "Set Kubemq gRPC port",
		}
		err = survey.AskOne(promptGrpcPort, &cfg.GrpcPort)
		if err != nil {
			return err
		}
		promptRestPort := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set Kubemq Rest port interface:",
			Default:  "9090",
			Help:     "Set Kubemq Rest port",
		}
		err = survey.AskOne(promptRestPort, &cfg.RestPort)
		if err != nil {
			return err
		}
		promptAPIPort := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set Kubemq Api port interface:",
			Default:  "8080",
			Help:     "Set Kubemq Api port",
		}
		err = survey.AskOne(promptAPIPort, &cfg.ApiPort)
		if err != nil {
			return err
		}
		connectionTypeSelect := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Set default interface:",
			Options:  []string{"grpc", "rest"},
			Default:  "grpc",
			Help:     "Select the default interface connection type",
		}
		err = survey.AskOne(connectionTypeSelect, &cfg.ConnectionType)
		if err != nil {
			return err
		}
	}
	return nil
}
