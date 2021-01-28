package config

import (
	"context"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"

	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"sort"
	"strings"
)

type ContextOptions struct {
	cfg *config.Config
}

var contextExamples = `
	# Execute context configuration
	# kubemqctl config context
`
var contextLong = `Config context command allows to set Kubemqctl context`
var contextShort = `Config context command allows to set Kubemqctl context`

// NewCmdCreate returns new initialized instance of create sub command
func NewCmdContext(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ContextOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{
		Use:     "context",
		Aliases: []string{"ctx", "c"},
		Short:   contextShort,
		Long:    contextLong,
		Example: contextExamples,
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

func (o *ContextOptions) Complete(args []string, cfg *config.Config) error {
	return nil
}

func (o *ContextOptions) Validate() error {
	return nil
}

func (o *ContextOptions) Run(ctx context.Context) error {
	err := runContextSelection(o.cfg)
	if err != nil {
		return err
	}
	return o.cfg.Save()
}

func runContextSelection(cfg *config.Config) error {
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

	list, err = getClusters(cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	if len(list) == 0 {
		utils.Println("No Kubemq clusters were found for selection")
	} else {
		clusterSelected := ""
		clusterSelect := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Select current Kubemq cluster:",
			Options:  list,
			Default:  list[0],
			Help:     "Select the default Kubemq cluster from available Kubemq clusters",
		}
		err := survey.AskOne(clusterSelect, &clusterSelected)
		if err != nil {
			return err
		}
		pair := strings.Split(clusterSelected, "/")
		cfg.CurrentNamespace = pair[0]
		cfg.CurrentStatefulSet = pair[1]
	}
	return nil
}
func getClusters(kubeConfig string) ([]string, error) {
	c, err := client.NewClient(kubeConfig)
	if err != nil {
		return nil, err
	}
	clusterManager, err := cluster.NewManager(c)
	if err != nil {
		return nil, err
	}
	clusters, err := clusterManager.GetKubemqClusters()
	if err != nil {
		return nil, err
	}

	return clusters.List(), err

}
