package cluster

import (
	"github.com/kubemq-io/kubetools/cmd/cluster/deploy"
	"github.com/kubemq-io/kubetools/cmd/cluster/logs"
	"github.com/kubemq-io/kubetools/cmd/cluster/proxy"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/spf13/cobra"
)

var clusterExamples = ``
var clusterLong = `Execute KubeMQ cluster commands`
var clusterShort = `Execute KubeMQ cluster commands`

// NewCmdCreate returns new initialized instance of create sub query
func NewCmdCluster(cfg *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"cl"},
		Short:   clusterLong,
		Long:    clusterShort,
		Example: clusterExamples,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(proxy.NewCmdProxy(cfg))
	cmd.AddCommand(logs.NewCmdLogs(cfg))
	cmd.AddCommand(deploy.NewCmdDeploy(cfg))

	return cmd
}
