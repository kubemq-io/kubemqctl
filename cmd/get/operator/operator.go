package operator

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type GetOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get Kubemq operators list 
	kubemqctl operator get  
`
var getLong = `Get command display all operators deployed across all namespaces`
var getShort = `Get Kubemq Operators List`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &GetOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "operator",
		Aliases: []string{"operators", "op", "o"},
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

	return cmd
}

func (o *GetOptions) Complete(args []string) error {
	return nil
}

func (o *GetOptions) Validate() error {
	return nil
}

func (o *GetOptions) Run(ctx context.Context) error {
	newClient, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	operatorManager, err := operator.NewManager(newClient)
	if err != nil {
		return err
	}

	utils.Println("Getting Kubemq Operators List...")
	operators, err := operatorManager.GetKubemqOperators()
	if err != nil {
		return err
	}
	if len(operators.List()) == 0 {
		return fmt.Errorf("no Kubemq operators were found in the cluster")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "OPERATOR\n")
	for _, item := range operators.List() {
		fmt.Fprintf(w, "%s\n",
			item,
		)
	}
	w.Flush()
	return nil
}
