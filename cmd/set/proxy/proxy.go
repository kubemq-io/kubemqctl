package proxy

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"

	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
)

type ProxyOptions struct {
	cfg *config.Config
	*k8s.ProxyOptions
}

var proxyExamples = `
	# Proxy a Kubemq cluster ports
	kubemqctl cluster proxy
`
var proxyLong = `Proxy command allows to act as a full layer 4 proxy (port-forwarding) of a Kubemq cluster connection to localhost. Proxy a KubeMW cluster allows the developer to interact with remote Kubemq cluster ports as localhost `
var proxyShort = `Proxy Kubemq cluster connection to localhost command`

func NewCmdProxy(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ProxyOptions{
		cfg: cfg,
		ProxyOptions: &k8s.ProxyOptions{
			KubeConfig: cfg.KubeConfigPath,
			Namespace:  "",
			Pod:        "",
			Ports:      nil,
		},
	}
	cmd := &cobra.Command{

		Use:     "proxy",
		Aliases: []string{"p"},
		Short:   proxyShort,
		Long:    proxyLong,
		Example: proxyExamples,
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

func (o *ProxyOptions) Complete(args []string) error {
	o.Ports = []string{"8080", "9090", "50000"}
	return nil
}

func (o *ProxyOptions) Validate() error {
	return nil
}

func (o *ProxyOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	list, err := c.GetKubemqClusters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no Kubemq clusters were found to proxy")
	}
	selection := ""
	multiSelected := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select Kubemq cluster to Proxy",
		Options:       list,
		Default:       list[0],
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(multiSelected, &selection)
	if err != nil {
		return err
	}
	ns, name := client.StringSplit(selection)
	o.Namespace = ns
	o.StatefulSet = name

	err = k8s.SetProxy(ctx, o.ProxyOptions)
	if err != nil {
		return err
	}
	return nil
}
