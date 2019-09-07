package proxy

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type ProxyOptions struct {
	cfg *config.Config
	*k8s.ProxyOptions
}

var proxyExamples = `
	# proxy default/kubemq-cluster with default KubeMQ ports
	kubetools cluster proxy

	# proxy specific namespace/pod with default KubeMQ ports
	kubetools cluster proxy kubemq kubemq-cluster1 

	# proxy specific namespace/pod with specific ports
	kubetools cluster proxy default nginx -p 80:80 
`
var proxyLong = `Proxy KubeMQ cluster connection to localhost`
var proxyShort = `Proxy KubeMQ cluster connection to localhost`

func NewCmdProxy(cfg *config.Config) *cobra.Command {
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
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringArrayVarP(&o.Ports, "ports", "p", []string{
		fmt.Sprintf("%d:%d", cfg.GrpcPort, cfg.GrpcPort),
		fmt.Sprintf("%d:%d", cfg.RestPort, cfg.RestPort),
		fmt.Sprintf("%d:%d", cfg.ApiPort, cfg.ApiPort),
	}, "Set proxy ports")

	return cmd
}

func (o *ProxyOptions) Complete(args []string) error {
	if len(args) >= 2 {
		o.Namespace = args[0]
		o.Pod = args[1]
		return nil
	} else {
		o.Namespace = o.cfg.CurrentNamespace
		o.StatefulSet = o.cfg.CurrentStatefulSet
	}
	return nil
}

func (o *ProxyOptions) Validate() error {
	return nil
}

func (o *ProxyOptions) Run(ctx context.Context) error {
	err := k8s.SetProxy(ctx, o.ProxyOptions)
	if err != nil {
		return err
	}
	return nil
}
