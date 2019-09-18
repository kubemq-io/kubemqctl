package apply

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemqctl/pkg/config"
	"github.com/kubemq-io/kubemqctl/pkg/k8s/deployment"
	"github.com/kubemq-io/kubemqctl/pkg/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
)

type ApplyOptions struct {
	cfg        *config.Config
	importData string
	fileName   string
	watch      bool
	status     bool
}

var applyExamples = `
	# Apply KubeMQ cluster deployment
	# kubemqctl cluster apply kubemq-cluster.yaml 

	# Apply KubeMQ cluster deployment with watching status and events
	# kubemqctl cluster apply kubemq-cluster.yaml -w -s

`
var applyLong = `Apply command allows an update to a KubeMQ StatefulSet configuration with a yaml file`
var applyShort = `Apply a KubeMQ cluster command`

func NewCmdApply(ctx context.Context, cfg *config.Config) *cobra.Command {
	o := &ApplyOptions{
		cfg: cfg,
	}
	cmd := &cobra.Command{

		Use:     "apply",
		Aliases: []string{"a", "ap"},
		Short:   applyShort,
		Long:    applyLong,
		Example: applyExamples,
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			utils.CheckErr(o.Complete(args), cmd)
			utils.CheckErr(o.Validate())
			utils.CheckErr(o.Run(ctx))
		},
	}
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "watch and print apply StatefulSet events")
	cmd.PersistentFlags().BoolVarP(&o.status, "status", "s", false, "watch and print apply StatefulSet status")
	cmd.PersistentFlags().StringVarP(&o.fileName, "file", "f", "", "set yaml configuration file")
	cmd.MarkFlagRequired("file")
	cmd.MarkFlagFilename("file", "yaml", "yml")
	return cmd
}

func (o *ApplyOptions) Complete(args []string) error {

	if o.fileName != "" {
		buff, err := ioutil.ReadFile(o.fileName)
		if err != nil {
			return err
		}
		o.importData = string(buff)
		return nil
	}
	return fmt.Errorf("invalid file name to import")
}

func (o *ApplyOptions) Validate() error {
	if o.importData == "" {
		return fmt.Errorf("configuration file empty")
	}
	return nil
}

func (o *ApplyOptions) Run(ctx context.Context) error {
	sd, err := deployment.NewStatefulSetDeployment(o.cfg)
	if err != nil {
		return err
	}
	err = sd.Import(o.importData)
	if err != nil {
		return err
	}
	utils.Printlnf("Apply started...")

	executed, err := sd.Execute(sd.StatefulSet.Name, sd.StatefulSet.Namespace)
	if err != nil {
		return err
	}
	if !executed {
		return nil
	}
	stsName := sd.StatefulSet.Name
	stsNamespace := sd.StatefulSet.Namespace

	if o.watch {
		go sd.Client.PrintEvents(ctx, stsNamespace, stsName)
	}

	if o.status {
		go sd.Client.PrintStatefulSetStatus(ctx, *sd.StatefulSet.Spec.Replicas, stsNamespace, stsName)
	}
	if o.status || o.watch {
		<-ctx.Done()

	}
	return nil
}
