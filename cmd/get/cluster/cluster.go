package cluster

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type getOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get status of Kubemq of clusters
	kubemqctl get clusters
`
var statusLong = `Get command allows to show the current information of requested resource`
var statusShort = `Get information of requested resource`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &getOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "cluster",
		Aliases: []string{"c", "clusters"},
		Short:   statusShort,
		Long:    statusLong,
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
	client, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	clusterManager, err := cluster.NewManager(client)
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
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tDESIRED\tREADY\tIMAGE\tGRPC\tREST\tAPI\n")
	for _, name := range clusters.List() {
		cluster := clusters.Cluster(name)

		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t%s\t%s\t%s\n",
			name,
			*cluster.Status.Replicas,
			cluster.Status.Ready,
			cluster.Status.Version,
			cluster.Status.Grpc,
			cluster.Status.Rest,
			cluster.Status.Api,
		)
	}
	w.Flush()
	return nil
}
