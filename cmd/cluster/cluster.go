package cluster

import (
	"github.com/kubemq-io/kubemqctl/cmd/cluster/apply"
	context "github.com/kubemq-io/kubemqctl/cmd/cluster/context"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/create"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/dashboard"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/delete"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/describe"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/list"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/logs"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/proxy"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/scale"
	"github.com/kubemq-io/kubemqctl/pkg/config"

	"github.com/spf13/cobra"
)

var clusterExamples = `
	# Execute create KubeMQ cluster command
	kubemqctl cluster create

	# Execute delete KubeMQ cluster command
	kubemqctl cluster delete

	# Execute describe KubeMQ cluster command
	kubemqctl cluster describe

	# Execute apply KubeMQ cluster command
	kubemqctl cluster apply

	# Execute show KubeMQ cluster logs command
	kubemqctl cluster logs

	# Execute scale KubeMQ cluster command
	kubemqctl cluster scale

	# Execute list of KubeMQ clusters command
	kubemqctl cluster list

	# Execute proxy ports of KubeMQ cluster command
	kubemqctl cluster proxy

	# Execute switch Kubernetes context command
	kubemqctl cluster context

	# Execute cluster web-based dashboard
	kubemqctl cluster dashboard
`
var clusterLong = `Executes KubeMQ cluster management commands`
var clusterShort = `Executes KubeMQ cluster management commands`

func NewCmdCluster(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "cluster",
		Aliases:   []string{"c", "cls"},
		Short:     clusterShort,
		Long:      clusterLong,
		Example:   clusterExamples,
		ValidArgs: []string{"create", "context", "apply", "dashboard", "delete", "describe", "list", "logs", "proxy", "scale"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(logs.NewCmdLogs(cfg))
	cmd.AddCommand(create.NewCmdCreate(cfg))
	cmd.AddCommand(delete.NewCmdDelete(cfg))
	cmd.AddCommand(scale.NewCmdScale(cfg))
	cmd.AddCommand(list.NewCmdList(cfg))
	cmd.AddCommand(context.NewCmdContext(cfg))
	cmd.AddCommand(proxy.NewCmdProxy(cfg))
	cmd.AddCommand(apply.NewCmdApply(cfg))
	cmd.AddCommand(describe.NewCmdDescribe(cfg))
	cmd.AddCommand(dashboard.NewCmdDashboard(cfg))
	return cmd
}
