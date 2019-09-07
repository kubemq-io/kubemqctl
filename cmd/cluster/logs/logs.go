package logs

import (
	"context"
	"github.com/AlecAivazis/survey/v2"
	"strings"

	"github.com/kubemq-io/kubetools/pkg/config"

	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/k8s/logs"

	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
)

type LogsOptions struct {
	cfg *config.Config
	*logs.Options
	disableColor bool
}

var logsExamples = `
	# Stream logs with selection of KubeMQ cluster
	kubetools cluster logs

	# Stream logs of all pods in default namespace
	kubetools cluster logs .* -n default

	# Stream logs of regex base pods with logs since 10m ago
	kubetools cluster logs kubemq-cluster.* -s 10m

	# Stream logs of regex base pods with logs since 10m ago include the string of 'connection'
	kubetools cluster logs kubemq-cluster.* -s 10m -i connection

	# Stream logs of regex base pods with logs exclude the string of 'error'
	kubetools cluster logs kubemq-cluster.* -s 10m -e error

	# Stream logs of specific container
	kubetools cluster logs -c kubemq-cluster-0
`
var logsLong = `Stream logs from KubeMQ pods`
var logsShort = `Stream logs from KubeMQ pods`

func NewCmdLogs(cfg *config.Config) *cobra.Command {
	o := &LogsOptions{
		cfg: cfg,
		Options: &logs.Options{
			PodQuery:       ".*",
			ContainerQuery: ".*",
			Timestamps:     false,
			Since:          0,
			Namespace:      "",
			Exclude:        nil,
			Include:        nil,
			AllNamespaces:  true,
			Selector:       "",
			Tail:           0,
			Color:          "auto",
		},
	}
	cmd := &cobra.Command{

		Use:     "logs",
		Aliases: []string{"lgs"},
		Short:   logsShort,
		Long:    logsLong,
		Example: logsExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().DurationVarP(&o.Options.Since, "since", "s", 0, "Set since duration time")
	cmd.PersistentFlags().StringVarP(&o.Options.Namespace, "namespace", "n", "", "Set default namespace")
	cmd.PersistentFlags().StringVarP(&o.Options.ContainerQuery, "container", "c", "", "Set container regex")
	cmd.PersistentFlags().StringArrayVarP(&o.Options.Include, "include", "i", []string{}, "Set strings to include")
	cmd.PersistentFlags().StringArrayVarP(&o.Options.Exclude, "exclude", "e", []string{}, "Set strings to exclude")
	cmd.PersistentFlags().StringVarP(&o.Options.Selector, "label", "l", "", "Set label selector")
	cmd.PersistentFlags().Int64VarP(&o.Options.Tail, "tail", "t", 0, "Set how many lines to tail for each pod")
	cmd.PersistentFlags().BoolVarP(&o.disableColor, "disable-color", "", false, "Set to disable colorized output")

	return cmd
}

func (o *LogsOptions) Complete(args []string) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		list, err := c.GetKubeMQClusters()
		if err != nil {
			return err
		}
		if len(list) == 0 {
			goto NEXT
		}
		selection := ""
		prompt := &survey.Select{
			Renderer: survey.Renderer{},
			Message:  "Show logs for KubeMQ cluster:",
			Options:  list,
			Default:  list[0],
		}
		err = survey.AskOne(prompt, &selection)
		if err != nil {
			return err
		}
		pair := strings.Split(selection, "/")
		o.Options.Namespace = pair[0]
		o.Options.PodQuery = pair[1]
	}
NEXT:
	if len(args) > 0 {
		o.PodQuery = args[0]
	}
	if o.Options.Namespace != "" {
		o.Options.AllNamespaces = false
	}
	if o.disableColor {
		o.Options.Color = "never"
	}
	return nil
}

func (o *LogsOptions) Validate() error {
	return nil
}

func (o *LogsOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	err = logs.Run(ctx, c, o.Options)
	if err != nil {
		return err
	}
	return nil
}
