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
	watch      bool
	status     bool
}

var applyExamples = `
	# Apply KubeMQ cluster deployment
	# kubemqctl cluster apply kubemq-cluster.yaml 

	# Apply KubeMQ cluster deployment with watching status and events
	# kubemqctl cluster apply kubemq-cluster.yaml -w -s

`
var applyLong = `Apply a KubeMQ cluster`
var applyShort = `Apply a KubeMQ cluster`

func NewCmdApply(cfg *config.Config) *cobra.Command {
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
	cmd.PersistentFlags().BoolVarP(&o.watch, "watch", "w", false, "watch and print apply statefulset events")
	cmd.PersistentFlags().BoolVarP(&o.status, "status", "s", false, "watch and print apply statefulset status")

	return cmd
}

func (o *ApplyOptions) Complete(args []string) error {
	if len(args) != 0 {
		buff, err := ioutil.ReadFile(args[0])
		if err != nil {
			return err
		}
		o.importData = string(buff)
		return nil
	}
	return fmt.Errorf("invalid argument, no yaml file specified")
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
