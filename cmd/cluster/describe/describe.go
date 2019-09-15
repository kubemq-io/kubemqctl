package describe

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

type DescribeOptions struct {
	cfg  *config.Config
	file bool
}

var describeExamples = `
 	# Describe KubeMQ cluster to console
	kubemqctl cluster describe

	# Describe KubeMQ cluster to a file
	kubemqctl cluster describe -f
`
var describeLong = `Describe KubeMQ cluster`
var describeShort = `Describe KubeMQ cluster`

func NewCmdDescribe(cfg *config.Config) *cobra.Command {
	o := &DescribeOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "describe",
		Aliases: []string{"des", "ds"},
		Short:   describeShort,
		Long:    describeLong,
		Example: describeExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().BoolVarP(&o.file, "file", "f", false, "export to yaml file")
	return cmd
}

func (o *DescribeOptions) Complete(args []string) error {
	return nil
}

func (o *DescribeOptions) Validate() error {

	return nil
}

func (o *DescribeOptions) Run(ctx context.Context) error {
	c, err := client.NewClient(o.cfg.KubeConfigPath)
	if err != nil {
		return err
	}

	list, err := c.GetKubeMQClusters()
	if err != nil {
		return err
	}
	if len(list) == 0 {
		return fmt.Errorf("no KubeMQ clusters were found to describe")
	}
	selection := ""
	multiSelected := &survey.Select{
		Renderer:      survey.Renderer{},
		Message:       "Select KubeMQ cluster to describe",
		Options:       list,
		Default:       list[0],
		Help:          "Select KubeMQ cluster to describe",
		PageSize:      0,
		VimMode:       false,
		FilterMessage: "",
		Filter:        nil,
	}
	err = survey.AskOne(multiSelected, &selection)
	if err != nil {
		return err
	}
	ns, name := client.StringSplit(selection)
	sd, err := deployment.NewStatefulSetDeploymentFromCluster(o.cfg, ns, name)

	if err != nil {
		return err
	}

	if o.file {
		f, err := os.Create(fmt.Sprintf("%s-%s.yaml", ns, name))
		if err != nil {
			return err
		}
		err = sd.Export(f)
		if err != nil {
			utils.Printlnf("export to file %s-%s.yaml failed", ns, name)
			return err
		}
		utils.Printlnf("export to file %s-%s.yaml completed", ns, name)

	} else {
		_ = sd.Export(os.Stdout)
	}

	return nil
}
