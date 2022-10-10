package dashboard

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"strings"
)

type getOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get KubeMQ web interface
	kubemqctl get dashboard
`
var getLong = `Get access to KubeMQ dashboard`
var getShort = `Get access to KubeMQ dashboard`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &getOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "dashboard",
		Aliases: []string{"d", "dash", "dashboard"},
		Short:   getShort,
		Long:    getLong,
		Example: getExamples,
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

func (o *getOptions) Complete(args []string) error {

	return nil
}

func (o *getOptions) Validate() error {

	return nil
}

func (o *getOptions) Run(ctx context.Context) error {
	namespace, clusterName := "", ""
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	clusterManager, err := cluster.NewManager(c)
	if err != nil {
		return err
	}

	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return err
	}

	if len(clusters.List()) == 0 {
		return fmt.Errorf("no Kubemq clusters were found")
	}
	if len(clusters.List()) == 1 {
		pair := clusters.List()[0]
		namespace, clusterName = StringSplit(pair)
	} else {
		selection := ""
		prompt := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Show Dashboard for KubeMQ cluster:",
			Options:  clusters.List(),
			Default:  clusters.List()[0],
		}
		err = survey.AskOne(prompt, &selection)
		if err != nil {
			return err
		}
		pair := strings.Split(selection, "/")
		namespace = pair[0]
		clusterName = pair[1]
	}

	fmt.Printf("Opening Dashboard for Kubemq cluster: %s/%s\n", namespace, clusterName)
	proxyOptions := &k8s.ProxyOptions{
		KubeConfig:  o.cfg.KubeConfigPath,
		Namespace:   namespace,
		StatefulSet: clusterName,
		Pod:         "",
		Ports:       []string{"8080"},
	}
	errCh := make(chan error, 1)
	err = k8s.SetConcurrentProxy(ctx, proxyOptions, errCh)
	if err != nil {
		return err
	}
	err = browser.OpenURL("http://localhost:8080")
	if err != nil {
		return err

	}
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return nil
	}

}

func StringSplit(input string) (string, string) {
	pair := strings.Split(input, "/")
	if len(pair) == 2 {
		return pair[0], pair[1]
	}
	return "", ""
}
