package config

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

type ConfigOptions struct {
	Cfg *config.Config
}

var configExamples = `
	# Run configuration wizard
	# kubetools config
`
var configLong = `config kubetools`
var configShort = `config kubetools`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdConfig(cfg *config.Config) *cobra.Command {
	o := &ConfigOptions{}
	cmd := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf"},
		Short:   configShort,
		Long:    configLong,

		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args, cfg))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
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
	cfg := o.Cfg
	integrationType := ""
	integrationSelect := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Select KubeMQ location:",
		Options:  []string{"kubernetes cluster", "single docker container"},
		Default:  "kubernetes cluster",
		Help:     "select the location of KubeMQ server",
	}
	err := survey.AskOne(integrationSelect, &integrationType)
	if err != nil {
		return err
	}
	if integrationType == "kubernetes cluster" {
		cfg.AutoIntegrated = true
		kubeConfigPath := ""
		prompt := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Select kube config path (press Enter for default):",
			Default:  "",
			Help:     "set kube.config file path if not kubectl default",
		}
		err := survey.AskOne(prompt, &kubeConfigPath)
		if err != nil {
			return err
		}
		if kubeConfigPath != "" {
			cfg.KubeConfigPath = kubeConfigPath
		}

		c, err := client.NewClient(kubeConfigPath)
		if err != nil {
			return err
		}
		contextMap, current, err := c.GetConfigContext()
		if err != nil {
			return err
		}
		list := []string{}
		for key, _ := range contextMap {
			list = append(list, key)
		}
		sort.Strings(list)
		contextSelected := ""
		contextSelect := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       "Select kubernetes cluster context",
			Options:       list,
			Default:       current,
			Help:          "set kubernetes connection context",
			PageSize:      0,
			VimMode:       false,
			FilterMessage: "",
			Filter:        nil,
		}
		err = survey.AskOne(contextSelect, &contextSelected)
		if err != nil {
			return err
		}
		err = c.SwitchContext(contextSelected)
		if err != nil {
			return err
		}

		list, err = o.getClusters(kubeConfigPath)
		if err != nil {
			return err
		}
		if list == nil {
			utils.Println("No KubeMQ clusters were found for selection")
		} else {
			clusterSelected := ""
			clusterSelect := &survey.Select{
				Renderer: survey.Renderer{},
				Message:  "Select current KubeMQ cluster:",
				Options:  list,
				Default:  list[0],
				Help:     "select the default KubeMQ cluster from available KubeMQ clusters ",
			}
			err := survey.AskOne(clusterSelect, &clusterSelected)
			if err != nil {
				return err
			}
			pair := strings.Split(clusterSelected, "/")
			cfg.CurrentNamespace = pair[0]
			cfg.CurrentStatefulSet = pair[1]
		}
	} else {
		cfg.AutoIntegrated = false
	}
	promptHost := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Host:",
		Default:  "localhost",
		Help:     "set KubeMQ host",
	}
	err = survey.AskOne(promptHost, &cfg.Host, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return err
	}
	promptGrpcPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ gRPC port interface:",
		Default:  "50000",
		Help:     "set KubeMQ gRPC port",
	}
	err = survey.AskOne(promptGrpcPort, &cfg.GrpcPort)
	if err != nil {
		return err
	}
	promptRestPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Rest port interface:",
		Default:  "9090",
		Help:     "set KubeMQ Rest port",
	}
	err = survey.AskOne(promptRestPort, &cfg.RestPort)
	if err != nil {
		return err
	}
	promptAPIPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Api port interface:",
		Default:  "8080",
		Help:     "set KubeMQ Api port",
	}
	err = survey.AskOne(promptAPIPort, &cfg.ApiPort)
	if err != nil {
		return err
	}
	connextionTypeSelect := &survey.Select{
		Renderer: survey.Renderer{},
		Message:  "Set default interface:",
		Options:  []string{"grpc", "rest"},
		Default:  "grpc",
		Help:     "select the default interface connection type",
	}
	err = survey.AskOne(connextionTypeSelect, &cfg.ConnectionType)
	if err != nil {
		return err
	}
	promptIsSecured := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Set secured connection ?:",
		Default:  false,
		Help:     "set KubeMQ secure connection",
	}
	err = survey.AskOne(promptIsSecured, &cfg.IsSecured)
	if err != nil {
		return err
	}
	if cfg.IsSecured {
		promptCertFile := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set cert file path:",
			Default:  "./",
			Help:     "set KubeMQ cert file path",
		}
		err = survey.AskOne(promptCertFile, &cfg.CertFile)
		if err != nil {
			return err
		}
	}
	err = cfg.Save()
	if err != nil {
		return err
	}
	utils.Println("configuration completed and saved.")
	return nil
}

func (o *ConfigOptions) getClusters(kubeConfig string) ([]string, error) {
	c, err := client.NewClient(kubeConfig)
	if err != nil {
		return nil, err
	}
	return c.GetKubeMQClusters()

}
