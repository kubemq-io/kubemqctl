package cluster

import (
	context "github.com/kubemq-io/kubetools/cmd/cluster/context"
	"github.com/kubemq-io/kubetools/cmd/cluster/create"
	"github.com/kubemq-io/kubetools/cmd/cluster/delete"
	"github.com/kubemq-io/kubetools/cmd/cluster/list"
	"github.com/kubemq-io/kubetools/cmd/cluster/logs"
	"github.com/kubemq-io/kubetools/cmd/cluster/proxy"
	"github.com/kubemq-io/kubetools/cmd/cluster/scale"
	"github.com/kubemq-io/kubetools/pkg/config"

	"github.com/spf13/cobra"
)

var clusterExamples = `

`
var clusterLong = `Executes KubeMQ cluster management commands`
var clusterShort = `Executes KubeMQ cluster management commands`

func NewCmdCluster(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{

		Use:     "cluster",
		Aliases: []string{"cls"},
		Short:   clusterShort,
		Long:    clusterLong,
		Example: clusterExamples,
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
	return cmd
}
