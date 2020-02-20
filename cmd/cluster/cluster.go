package cluster

import (
	"context"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/apply"

	kubeCtx "github.com/kubemq-io/kubemqctl/cmd/cluster/context"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/create"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/dashboard"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/delete"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/describe"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/events"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/get"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/logs"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/proxy"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/register"
	"github.com/kubemq-io/kubemqctl/cmd/cluster/scale"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/utils"

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

	# Execute get of KubeMQ clusters command
	kubemqctl cluster get

	# Execute proxy ports of KubeMQ cluster command
	kubemqctl cluster proxy

	# Execute switch Kubernetes context command
	kubemqctl cluster context

	# Execute cluster web-based dashboard
	kubemqctl cluster dashboard

	# Show cluster events
	kubemqctl cluster events
`
var clusterLong = `Executes KubeMQ cluster management commands`
var clusterShort = `Executes KubeMQ cluster management commands`

func NewCmdCluster(ctx context.Context, cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:       "cluster",
		Aliases:   []string{"c", "cls"},
		Short:     clusterShort,
		Long:      clusterLong,
		Example:   clusterExamples,
		ValidArgs: []string{"create", "context", "apply", "dashboard", "delete", "describe", "get", "logs", "proxy", "scale", "events"},
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckErr(cmd.Help())
		},
	}
	cmd.AddCommand(logs.NewCmdLogs(ctx, cfg))
	cmd.AddCommand(create.NewCmdCreate(ctx, cfg))
	cmd.AddCommand(delete.NewCmdDelete(ctx, cfg))
	cmd.AddCommand(scale.NewCmdScale(ctx, cfg))
	cmd.AddCommand(get.NewCmdList(ctx, cfg))
	cmd.AddCommand(kubeCtx.NewCmdContext(ctx, cfg))
	cmd.AddCommand(proxy.NewCmdProxy(ctx, cfg))
	cmd.AddCommand(apply.NewCmdApply(ctx, cfg))
	cmd.AddCommand(describe.NewCmdDescribe(ctx, cfg))
	cmd.AddCommand(dashboard.NewCmdDashboard(ctx, cfg))
	cmd.AddCommand(events.NewCmdEvents(ctx, cfg))
	cmd.AddCommand(register.NewCmdRegister(ctx, cfg))
	return cmd
}
