package status

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubetools/pkg/config"
	"github.com/kubemq-io/kubetools/pkg/k8s/client"
	"github.com/kubemq-io/kubetools/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
	"time"
)

type StatusOptions struct {
	cfg *config.Config
}

var statusExamples = `
	# Get status of KubeMQ of clusters
	kubetools status
`
var statusLong = `Get status of KubeMQ of clusters`
var statusShort = `Get status of KubeMQ of clusters`

func NewCmdStatus(cfg *config.Config) *cobra.Command {
	o := &StatusOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "status",
		Aliases: []string{"st"},
		Short:   statusShort,
		Long:    statusLong,
		Example: statusExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}

	return cmd
}

func (o *StatusOptions) Complete(args []string) error {
	return nil
}

func (o *StatusOptions) Validate() error {

	return nil
}

func (o *StatusOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	list, err := c.GetKubeMQClustersStatus()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "NAME\tDESIRED\tRUNNING\tREADY\tIMAGE\tAGE\tSERVICES\n")
	for _, item := range list {
		dep, err := c.GetStatefulSetDeployment(item.Namespace, item.Name)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "%s/%s\t%d\t%d\t%d\t%s\t%s\t%s\n",
			dep.StatefulSetStatus.Namespace,
			dep.StatefulSetStatus.Name,
			dep.StatefulSetStatus.Desired,
			dep.StatefulSetStatus.Running,
			dep.StatefulSetStatus.Ready,
			dep.StatefulSetStatus.Image,
			dep.StatefulSetStatus.Age.Round(time.Second),
			dep.ServicesStatusString())
	}
	w.Flush()
	return nil
}
