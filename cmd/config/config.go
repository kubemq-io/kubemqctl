package config

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"sort"
	"strings"
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

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdConfig(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ConfigOptions{}
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
	cmd.AddCommand(NewCmdGetLicense(ctx))
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
		Message:  "Select KubeMQ install location:",
		Options: []string{
			"Kubernetes cluster",
			"MicroK8s",
			"K3s",
			"Minikube",
			"Other Kubernetes distribution",
			"Local docker container"},
		Default: "Kubernetes cluster",
		Help:    "Select the location of KubeMQ server",
	}
	err := survey.AskOne(integrationSelect, &integrationType)
	if err != nil {
		return err
	}
	if integrationType == "Local docker container" {
		cfg.AutoIntegrated = false
	} else {
		cfg.AutoIntegrated = true
		switch integrationType {
		case "Kubernetes cluster", "Minikube":
			cfg.KubeConfigPath = ""
		case "MicroK8s":
			cfg.KubeConfigPath = "/var/snap/microk8s/current/credentials/client.config"
		case "K3s":
			cfg.KubeConfigPath = "/etc/rancher/k3s/k3s.yaml"
		case "Other Kubernetes distribution":
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
		}
		c, err := client.NewClient(cfg.KubeConfigPath)
		if err != nil {
			return err
		}
		contextMap, current, err := c.GetConfigContext()
		if err != nil {
			return err
		}
		list := []string{}
		for key := range contextMap {
			list = append(list, key)
		}
		sort.Strings(list)
		contextSelected := ""
		contextSelect := &survey.Select{
			Renderer:      survey.Renderer{},
			Message:       "Select kubernetes cluster context:",
			Options:       list,
			Default:       current,
			Help:          "Set kubernetes connection context",
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

		list, err = o.getClusters(cfg.KubeConfigPath)
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
				Help:     "Select the default KubeMQ cluster from available KubeMQ clusters",
			}
			err := survey.AskOne(clusterSelect, &clusterSelected)
			if err != nil {
				return err
			}
			pair := strings.Split(clusterSelected, "/")
			cfg.CurrentNamespace = pair[0]
			cfg.CurrentStatefulSet = pair[1]
		}
	}

	promptHost := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Host:",
		Default:  "localhost",
		Help:     "Set KubeMQ host",
	}
	err = survey.AskOne(promptHost, &cfg.Host, survey.WithValidator(survey.MinLength(1)))
	if err != nil {
		return err
	}
	promptGrpcPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ gRPC port interface:",
		Default:  "50000",
		Help:     "Set KubeMQ gRPC port",
	}
	err = survey.AskOne(promptGrpcPort, &cfg.GrpcPort)
	if err != nil {
		return err
	}
	promptRestPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Rest port interface:",
		Default:  "9090",
		Help:     "Set KubeMQ Rest port",
	}
	err = survey.AskOne(promptRestPort, &cfg.RestPort)
	if err != nil {
		return err
	}
	promptAPIPort := &survey.Input{
		Renderer: survey.Renderer{},
		Message:  "Set KubeMQ Api port interface:",
		Default:  "8080",
		Help:     "Set KubeMQ Api port",
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
	promptIsSecured := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Set secured connection ?:",
		Default:  false,
		Help:     "Set KubeMQ secured connection",
	}
	err = survey.AskOne(promptIsSecured, &cfg.IsSecured)
	if err != nil {
		return err
	}
	if cfg.IsSecured {
		promptCertFile := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set cert file path:",
			Default:  "",
			Help:     "Set KubeMQ cert file path",
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
		Message:  "Would you like to set connection ClientId ?:",
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
		Message:  "Would you like to set JWT Authentication token ?:",
		Default:  false,
	}
	err = survey.AskOne(promptSetAuthToken, &isAuthToken)
	if err != nil {
		return err
	}
	if isAuthToken {
		promptAuthToken := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set JWT Authentication token file",
			Default:  cfg.AuthTokenFile,
			Help:     "Set JWT Authentication token file",
		}
		err = survey.AskOne(promptAuthToken, &cfg.AuthTokenFile)
		if err != nil {
			return err
		}
		_, err = ioutil.ReadFile(cfg.AuthTokenFile)
		if err != nil {
			return fmt.Errorf("error loading Authentication token file: %s", err.Error())
		}
	} else {
		cfg.AuthTokenFile = ""
	}
	isLicenseData := false
	promptSetLicenseData := &survey.Confirm{
		Renderer: survey.Renderer{},
		Message:  "Would you like to set KubeMQ enterprise license file ?:",
		Default:  false,
	}
	err = survey.AskOne(promptSetLicenseData, &isLicenseData)
	if err != nil {
		return err
	}
	if isLicenseData {
		dataFile := "license.key"
		promptLicenseDataFile := &survey.Input{
			Renderer: survey.Renderer{},
			Message:  "Set KubeMQ enterprise license file",
			Default:  dataFile,
			Help:     "Set KubeMQ enterprise license file",
		}
		err = survey.AskOne(promptLicenseDataFile, &dataFile)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadFile(dataFile)
		if err != nil {
			return fmt.Errorf("error loading license data file: %s", err.Error())
		}
		cfg.LicenseData = string(data)
	} else {
		cfg.LicenseData = ""
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
