package get

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

type GetOptions struct {
	cfg *config.Config
}

var getExamples = `
	# Get KubeMQ operators list 
	kubemqctl operator get  
`
var getLong = `Get command display all operators deployed across all namespaces`
var getShort = `Get KubeMQ Operators List`

func NewCmdGet(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &GetOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "get",
		Aliases: []string{"g"},
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
	mng, err := manager.NewManager(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}
	utils.Println("Getting KubeMQ Operators List...")
	list, err := mng.GetOperatorList()
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintf(w, "NAMESPACE\n")
	for _, item := range list {
		fmt.Fprintf(w, "%s\n",
			item,
		)
	}
	w.Flush()
	return nil
}
