package dashboard

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/cmd/get/dashboard/describe"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/dashboard"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type getOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get status of Kubemq of dashboards
	kubemqctl get dashboards
`
var getLong = `Get information of Kubemq dashboard resources`
var getShort = `Get information of Kubemq dashboard resources`

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
	cmd.AddCommand(describe.NewCmdDescribe(ctx, cfg))
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
	dashboardManager, err := dashboard.NewManager(client)
	if err != nil {
		return err
	}

	dashboards, err := dashboardManager.GetKubemqDashboardes()
	if err != nil {
		return err
	}
	if len(dashboards.List()) == 0 {
		return fmt.Errorf("no Kubemq dashboards were found")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tSTATUS\tVIEW-ADDRESS\tPROMETHEUS-VERSION\tGRPC\tGRAFANA-VERSION\n")
	for _, name := range dashboards.List() {
		dashboard := dashboards.Dashboard(name)

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			name,
			dashboard.Status.Status,
			dashboard.Status.Address,
			dashboard.Status.PrometheusVersion,
			dashboard.Status.GrafanaVersion,
		)
	}
	w.Flush()
	return nil
}
