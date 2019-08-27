package proxy

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	prx "github.com/kubemq-io/kubetools/pkg/k8s/proxy"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type ProxyOptions struct {
	cfg *config.Config
	*prx.ProxyOptions
}

var proxyExamples = `
	# proxy default/kubemq-cluster-0 with default KubeMQ ports
	kubetools proxy

	# proxy specific namespace/pod with default KubeMQ ports
	kubetools proxy kubemq kubemq-cluster1-0 

	# proxy specific namespace/pod with specific ports
	kubetools proxy default nginx -p 80:80 
`
var proxyLong = `proxy namespace/pod with ports`
var proxyShort = `proxy namespace/pod with ports`

func NewCmdProxy(cfg *config.Config) *cobra.Command {
	o := &ProxyOptions{
		cfg: cfg,
		ProxyOptions: &prx.ProxyOptions{
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
			utils.CheckErr(o.Complete(args))
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().StringArrayVarP(&o.Ports, "ports", "p", []string{
		fmt.Sprintf("%d:%d", cfg.GrpcPort, cfg.GrpcPort),
		fmt.Sprintf("%d:%d", cfg.RestPort, cfg.RestPort),
		fmt.Sprintf("%d:%d", cfg.ApiPort, cfg.ApiPort),
	}, "set proxy ports")

	return cmd
}

func (o *ProxyOptions) Complete(args []string) error {
	if len(args) == 0 {
		o.Namespace = o.cfg.CurrentNamespace
		o.Pod = o.cfg.CurrentStatefulSet + "-0"
		return nil
	}
	if len(args) == 1 {
		o.Namespace = args[0]
		o.Pod = o.cfg.CurrentStatefulSet + "-0"
		return nil
	}
	if len(args) >= 2 {
		o.Namespace = args[0]
		o.Pod = args[1]
		return nil
	}
	return nil
}

func (o *ProxyOptions) Validate() error {
	return nil
}

func (o *ProxyOptions) Run(ctx context.Context) error {
	utils.Printf("Set Proxy on %s/%s with Ports: %s\n", o.Namespace, o.Pod, o.Ports)
	err := prx.SetProxy(ctx, o.ProxyOptions)
	if err != nil {
		return err
	}
	return nil
}
