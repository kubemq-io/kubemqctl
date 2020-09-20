package connector

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/cmd/get/connector/describe"
	"github.com/kubemq-io/kubemqctl/cmd/get/connector/logs"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"text/tabwriter"
)

type getOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get status of Kubemq of connectors
	kubemqctl get connectors
`
var getLong = `Get information of Kubemq connectors resources`
var getShort = `Get information of Kubemq connectors resources`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &getOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "connectors",
		Aliases: []string{"cn", "con", "connector"},
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
	cmd.AddCommand(logs.NewCmdLogs(ctx, cfg))
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
	connectorManager, err := connector.NewManager(client)
	if err != nil {
		return err
	}

	connectors, err := connectorManager.GetKubemqConnectors()
	if err != nil {
		return err
	}
	if len(connectors.List()) == 0 {
		return fmt.Errorf("no Kubemq connectors were found")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tNAMESPACE\tREPLICAS\tTYPE\tIMAGE\tAPI\tSTATUS\t\n")
	for _, name := range connectors.List() {
		connector := connectors.Connector(name)
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			name,
			connector.Namespace,
			connector.Status.Replicas,
			connector.Status.Type,
			connector.Status.Image,
			connector.Status.Api,
			connector.Status.Status,
		)
	}
	w.Flush()
	return nil
}

func StringSplit(input string) (string, string) {
	pair := strings.Split(input, "/")
	if len(pair) == 2 {
		return pair[0], pair[1]
	}
	return "", ""
}
